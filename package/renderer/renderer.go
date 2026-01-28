// Package renderer
package renderer

import (
	//"fmt"
	"sort"
	"unsafe"

	"test2/package/framebuffer"
	"test2/package/layer"
)

type Renderer struct {
	W     int
	H     int
	L     []*layer.Layer
	Dirty []layer.Rect
	Front *framebuffer.Framebuffer
	Back  *framebuffer.Framebuffer
}

func (r *Renderer) AddLayer(fb *framebuffer.Framebuffer, z int, blend layer.BlendMode) {
	l := &layer.Layer{
		FB:    fb,
		Z:     z,
		Blend: blend,
		// Dirty: make([]layer.Rect, 0),
	}
	r.L = append(r.L, l)
	sort.Slice(r.L, func(i, j int) bool {
		return r.L[i].Z < r.L[j].Z
	})
}

func (r *Renderer) MarkDirty(rect layer.Rect) {
	r.Dirty = append(r.Dirty, rect)
}

func (r *Renderer) ClearDirty() {
	r.Dirty = r.Dirty[:0]
}

func copyRect(dst *framebuffer.Framebuffer, src *framebuffer.Framebuffer, r layer.Rect) {
	copyW := min(r.W, dst.W-r.X, src.W-r.X)
	copyH := min(r.H, dst.H-r.Y, src.H-r.Y)
	for y := range copyH {
		srcOneRow := unsafe.Slice(&src.Buf[((r.Y+y)*src.W)+r.X], copyW)
		dstOneRow := unsafe.Slice(&dst.Buf[((r.Y+y)*dst.W)+r.X], copyW)
		copy(dstOneRow, srcOneRow)
	}
}

func blendCopyRect(dst *framebuffer.Framebuffer, src *framebuffer.Framebuffer, r layer.Rect, blend layer.BlendMode) {
	copyW := min(r.W, dst.W-r.X, src.W-r.X)
	copyH := min(r.H, dst.H-r.Y, src.H-r.Y)
	for y := range copyH {
		srcOneRow := unsafe.Slice(&src.Buf[((r.Y+y)*src.W)+r.X], copyW)
		dstOneRow := unsafe.Slice(&dst.Buf[((r.Y+y)*dst.W)+r.X], copyW)
		switch blend {
		case layer.BlendCopy:
			copy(dstOneRow, srcOneRow)
		case layer.BlendOr:
			for x := range copyW {
				if srcOneRow[x] != ' ' && srcOneRow[x] != 0 {
					dstOneRow[x] = srcOneRow[x]
				}
			}
		}
	}
}

func (r *Renderer) Render() {
	// back <- Front
	for _, d := range r.Dirty {
		// fmt.Printf("(%d)%d,%d,%d,%d\r\n", i, d.X, d.Y, d.W, d.H)
		copyRect(r.Back, r.Front, d)
	}
	// copy Layers
	for _, l := range r.L {
		for _, d := range r.Dirty {
			blendCopyRect(r.Back, l.FB, d, l.Blend)
		}
	}
	// swap
	r.Back, r.Front = r.Front, r.Back
	// ClearDirty
	r.ClearDirty()
}

func (r Renderer) OutputFront() *framebuffer.Framebuffer {
	return r.Front
}

func New(w, h int) *Renderer {
	front := framebuffer.New(w, h)
	back := framebuffer.New(w, h)
	return &Renderer{
		W:     w,
		H:     h,
		Front: front,
		Back:  back,
	}
}
