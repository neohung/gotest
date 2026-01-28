// Package layer
package layer

import "test2/package/framebuffer"

type BlendMode int

const (
	BlendCopy BlendMode = 1
	BlendOr   BlendMode = 2
)

type Rect struct {
	X, Y int
	W, H int
}

type Layer struct {
	FB *framebuffer.Framebuffer
	Z  int
	// Dirty []Rect
	Blend BlendMode
}

// func (l *Layer) MarkDirty(r Rect) {
// 	l.Dirty = append(l.Dirty, r)
// }

// func (l *Layer) CleanDirty(r Rect) {
// 	l.Dirty = l.Dirty[:0]
// }
