// Package player
package player

import (
	"test2/package/framebuffer"
	"test2/package/world"
)

type Player struct {
	X, Y int
	CH   rune
}

func (p *Player) RenderPlayer(cam *world.Camera, fb *framebuffer.Framebuffer) {
	fb.Clear(' ')

	sx := p.X - cam.X
	sy := p.Y - cam.Y

	if sx >= 0 && sx < cam.W && sy >= 0 && sy < cam.H {
		fb.View[sy][sx] = p.CH
	}
}
