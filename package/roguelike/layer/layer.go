// Package layer
package layer

import (
	"test/package/roguelike/framebuffer"
)

type BlendMode int

const (
	BlendCopy BlendMode = iota // 直接覆蓋
	BlendOr                    // 透明overlay (' ' 跳過)
)

type Rect struct {
	X, Y int
	W, H int
}

type Layer struct {
	FB    *framebuffer.Framebuffer
	Z     int
	Blend BlendMode
	Dirty []Rect
}

func (l *Layer) MarkDirty(rect Rect) {
	l.Dirty = append(l.Dirty, rect)
}
