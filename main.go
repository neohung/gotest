// Package main
package main

import (
	//"fmt"
	"fmt"
	"time"

	"test2/package/framebuffer"
	"test2/package/layer"
	"test2/package/player"
	"test2/package/renderer"
	"test2/package/world"

	"github.com/gdamore/tcell/v2"
)

func PresentFB(s tcell.Screen, fb *framebuffer.Framebuffer) {
	s.Clear()
	for y := range fb.H {
		for x := range fb.W {
			s.SetContent(x, y, fb.View[y][x], nil, tcell.StyleDefault)
		}
	}
	s.Show()
}

func DrawHUDLine(fb *framebuffer.Framebuffer, y int, msg string) {
	fb.Clear(' ') // HUD 每 frame 重畫（正確）

	x := 0
	for _, ch := range msg {
		if x >= fb.W {
			break
		}
		fb.View[y][x] = ch
		x++
	}
}

func main() {
	camW, camH := 40, 20
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	w, h := screen.Size()
	// fmt.Printf("%d,%d\r\n", w, h)
	myWorld := world.New(100, 100, world.Camera{10, 10, camW, camH})
	myWorld.GenerateMap()
	px, py := myWorld.FindSpawn()
	mapFB := framebuffer.New(camW, camH)

	actorFB := framebuffer.New(camW, camH)
	actorFB.Clear(' ')
	p := player.Player{X: px, Y: py, CH: '@'}
	myWorld.Cam.Follow(p.X, p.Y, myWorld.W, myWorld.H)
	myWorld.RenderToMapFB(mapFB)
	p.RenderPlayer(&myWorld.Cam, actorFB)

	hudFB := framebuffer.New(w, h-2)

	r := renderer.New(w, h)
	r.AddLayer(mapFB, 1, layer.BlendCopy)
	r.AddLayer(actorFB, 2, layer.BlendOr)
	r.AddLayer(hudFB, 3, layer.BlendOr)
	// Mark once for first init
	r.MarkDirty(layer.Rect{X: 0, Y: 0, W: w, H: h})

	defer screen.Fini()
	screen.Clear()
	// screen.EnableMouse()
	eventChan := make(chan tcell.Event)
	go func() {
		for {
			eventChan <- screen.PollEvent()
		}
	}()
loop:
	for {
		oldCamX, oldCamY := myWorld.Cam.X, myWorld.Cam.Y
		// var oldX, oldY int
		select {
		case ev := <-eventChan:
			if ev == nil {
				break loop
			}
			switch ev := ev.(type) {
			case *tcell.EventKey:
				// oldX = p.X
				// oldY = p.Y
				switch ev.Key() {
				case tcell.KeyEsc:
					break loop
				case tcell.KeyEnter:
					r.ClearDirty()
					// time.Sleep(3000 * time.Millisecond)
				case tcell.KeyUp:
					p.Y--
				case tcell.KeyDown:
					p.Y++
				case tcell.KeyLeft:
					p.X--
				case tcell.KeyRight:
					p.X++
				}
			}
		default:
			// fps60 task here
			myWorld.Cam.Follow(p.X, p.Y, myWorld.W, myWorld.H)
			if myWorld.Cam.X != oldCamX || myWorld.Cam.Y != oldCamY {
				myWorld.RenderToMapFB(mapFB)
				r.MarkDirty(layer.Rect{X: 0, Y: 0, W: myWorld.Cam.W, H: myWorld.Cam.H})
				// DrawHUDLine(hudFB, hudFB.H-1, fmt.Sprintf("%d,%d,%d,%d,%d,%d", p.X, p.Y, myWorld.Cam.X, myWorld.Cam.Y, oldCamX, oldCamY))
				// r.MarkDirty(layer.Rect{X: 0, Y: hudFB.H - 1, W: hudFB.W, H: 1})
			}
			DrawHUDLine(hudFB, hudFB.H-1, fmt.Sprintf("%d,%d,%d,%d,%d,%d", p.X, p.Y, myWorld.Cam.X, myWorld.Cam.Y, oldCamX, oldCamY))
			r.MarkDirty(layer.Rect{X: 0, Y: hudFB.H - 1, W: hudFB.W, H: 1})
			r.Render()
			PresentFB(screen, r.OutputFront())
		}
		time.Sleep(16 * time.Millisecond)
	}
}
