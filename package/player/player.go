// Package player
package player

import (
	"test2/package/framebuffer"
	"test2/package/world"
)

type Player struct {
	X, Y int
	CH   rune
	FOV  int
}

func (p *Player) CanMove(newX, newY int, w *world.World) {
	if newX < 0 || newX > w.W-1 || newY < 0 || newY > w.H-1 {
		return
	}
	if w.Tiles[newY][newX] == rune(world.Wall) {
		return
	}
	p.X, p.Y = newX, newY
}

func (p *Player) RenderPlayer(cam *world.Camera, fb *framebuffer.Framebuffer) {
	fb.Clear(' ')

	sx := p.X - cam.X
	sy := p.Y - cam.Y

	if sx >= 0 && sx < cam.W && sy >= 0 && sy < cam.H {
		fb.View[sy][sx] = p.CH
	}
}
