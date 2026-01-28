// Package main
package main

import (
	"time"

	"test2/package/framebuffer"
	"test2/package/layer"
	"test2/package/renderer"

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

type Player struct {
	X, Y int
	CH   rune
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	w, h := screen.Size()
	// fmt.Printf("%d,%d\r\n", w, h)
	f1 := framebuffer.New(20, h-2)
	f1.Clear('.')

	f2 := framebuffer.New(20, h-2)
	f2.Clear(' ')
	p := Player{5, 5, '@'}
	f2.View[p.Y][p.X] = p.CH
	r := renderer.New(w, h)
	r.AddLayer(f1, 1, layer.BlendCopy)
	r.AddLayer(f2, 2, layer.BlendCopy)
	// Mark once for first init
	r.MarkDirty(layer.Rect{0, 0, w, h})
	r.Render()
	r.MarkDirty(layer.Rect{0, 0, w, h})
	r.Render()
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
		var oldX, oldY int
		select {
		case ev := <-eventChan:
			if ev == nil {
				break loop
			}
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEsc:
					break loop
				case tcell.KeyEnter:
					// r.MarkDirty(layer.Rect{0, 0, w, h})
					// screen.Clear()
					// PresentFB(screen, r.Front)
					// time.Sleep(3000 * time.Millisecond)
					screen.Clear()
					// PresentFB(screen, r.Back)
					PresentFB(screen, f2)
					time.Sleep(3000 * time.Millisecond)
				case tcell.KeyUp:
					oldX = p.X
					oldY = p.Y
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{oldX, oldY, 1, 1})
					p.Y--
					f2.View[p.Y][p.X] = p.CH
					r.MarkDirty(layer.Rect{p.X, p.Y, 1, 1})
					r.Render()
				case tcell.KeyDown:
					oldX = p.X
					oldY = p.Y
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{oldX, oldY, 1, 1})
					p.Y++
					r.MarkDirty(layer.Rect{p.X, p.Y, 1, 1})
					f2.View[p.Y][p.X] = p.CH
					r.Render()
				case tcell.KeyLeft:
					oldX = p.X
					oldY = p.Y
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{oldX, oldY, 1, 1})
					p.X--
					r.MarkDirty(layer.Rect{p.X, p.Y, 1, 1})
					f2.View[p.Y][p.X] = p.CH
					r.Render()
				case tcell.KeyRight:
					oldX = p.X
					oldY = p.Y
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{oldX, oldY, 1, 1})
					p.X++
					f2.View[p.Y][p.X] = p.CH
					r.MarkDirty(layer.Rect{p.X, p.Y, 1, 1})
					r.Render()
				}
			}
		default:
			// fps60 task here
			// fmt.Print("aaa")
			PresentFB(screen, r.OutputFront())
		}
		time.Sleep(100 * time.Millisecond)
	}
}
