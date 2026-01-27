// Package renderer
package renderer

import (
	"sort"
	"unsafe"

	"test/package/roguelike/framebuffer"
	"test/package/roguelike/layer"
)

type Renderer struct {
	W, H   int
	layers []*layer.Layer
	front  *framebuffer.Framebuffer
	back   *framebuffer.Framebuffer
	dirty  []layer.Rect
}

func NewRenderer(w, h int) *Renderer {
	return &Renderer{
		W:     w,
		H:     h,
		front: framebuffer.New(w, h),
		back:  framebuffer.New(w, h),
	}
}

func (r *Renderer) MarkDirty(rect layer.Rect) {
	r.dirty = append(r.dirty, rect)
}

func copyRect(dst, src *framebuffer.Framebuffer, r layer.Rect) {
	copyW := min(r.W, dst.W-r.X, src.W-r.X)
	copyH := min(r.H, dst.H-r.Y, src.H-r.Y)
	// for y := 0; y < copyH; y++ {
	for y := range copyH {
		srcOneRow := unsafe.Slice(&src.Buf[(y+r.Y)*src.W+r.X], copyW)
		dstOneRow := unsafe.Slice(&dst.Buf[(y+r.Y)*dst.W+r.X], copyW)
		copy(dstOneRow, srcOneRow)
	}
}

// 依據l裡的blend值決定l.Buf對dst的rect區塊做copy或or
func blendRect(l *layer.Layer, dst *framebuffer.Framebuffer, r layer.Rect) {
	src := l.FB
	copyW := min(r.W, dst.W-r.X, src.W-r.X)
	copyH := min(r.H, dst.H-r.Y, src.H-r.Y)

	for y := range copyH {
		srcOneRow := unsafe.Slice(&src.Buf[(y+r.Y)*src.W+r.X], copyW)
		dstOneRow := unsafe.Slice(&dst.Buf[(y+r.Y)*dst.W+r.X], copyW)
		switch l.Blend {
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

func (r *Renderer) AddLayer(fb *framebuffer.Framebuffer, z int, blend layer.BlendMode) *layer.Layer {
	l := &layer.Layer{
		FB:    fb,
		Z:     z,
		Blend: blend,
	}
	r.layers = append(r.layers, l)
	sort.Slice(r.layers, func(i, j int) bool {
		return r.layers[i].Z < r.layers[j].Z
	})
	return l
}

func (r *Renderer) Render() {
	// 1. back <- front 只copy MarkDirty, 先將front的內容複製到back
	for _, d := range r.dirty {
		copyRect(r.back, r.front, d)
	}
	// 2. Layers (Z-order)
	for _, l := range r.layers {
		// 合併 layer 的 dirty
		for _, d := range r.dirty {
			// r.layers的每個layer對r.back做區塊的blend操作, 等於每層l對back做區域render
			blendRect(l, r.back, d)
		}
	}
	// 3. swap
	r.front, r.back = r.back, r.front
	// 4. clear dirty
	r.dirty = r.dirty[:0]
}

// Front Present到 screen
func (r *Renderer) Front() *framebuffer.Framebuffer {
	return r.front
}
