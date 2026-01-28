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
	f1 := framebuffer.New(w, h-2)
	f1.Clear('.')

	f2 := framebuffer.New(w, h-2)
	f2.Clear(' ')
	p := Player{w / 2, h / 2, '@'}
	f2.View[p.Y][p.X] = p.CH
	r := renderer.New(w, h)
	r.AddLayer(f1, 1, layer.BlendCopy)
	r.AddLayer(f2, 2, layer.BlendOr)
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
		var oldX, oldY int
		select {
		case ev := <-eventChan:
			if ev == nil {
				break loop
			}
			switch ev := ev.(type) {
			case *tcell.EventKey:
				oldX = p.X
				oldY = p.Y
				switch ev.Key() {
				case tcell.KeyEsc:
					break loop
				case tcell.KeyEnter:
					r.ClearDirty()
					// time.Sleep(3000 * time.Millisecond)
				case tcell.KeyUp:
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{X: oldX, Y: oldY, W: 1, H: 1})
					p.Y--
					f2.View[p.Y][p.X] = p.CH
					r.MarkDirty(layer.Rect{X: p.X, Y: p.Y, W: 1, H: 1})
				case tcell.KeyDown:
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{X: oldX, Y: oldY, W: 1, H: 1})
					p.Y++
					r.MarkDirty(layer.Rect{X: p.X, Y: p.Y, W: 1, H: 1})
					f2.View[p.Y][p.X] = p.CH
				case tcell.KeyLeft:
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{X: oldX, Y: oldY, W: 1, H: 1})
					p.X--
					r.MarkDirty(layer.Rect{X: p.X, Y: p.Y, W: 1, H: 1})
					f2.View[p.Y][p.X] = p.CH
				case tcell.KeyRight:
					// f2.Clear(' ')
					f2.View[oldY][oldX] = ' '
					r.MarkDirty(layer.Rect{X: oldX, Y: oldY, W: 1, H: 1})
					p.X++
					f2.View[p.Y][p.X] = p.CH
					r.MarkDirty(layer.Rect{X: p.X, Y: p.Y, W: 1, H: 1})
				}
			}
		default:
			// fps60 task here
			r.Render()
			PresentFB(screen, r.OutputFront())
		}
		time.Sleep(16 * time.Millisecond)
	}
}
