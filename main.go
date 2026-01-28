// Package main
package main

import (
	"fmt"

	"test2/package/framebuffer"
	"test2/package/layer"
	"test2/package/renderer"

	"github.com/gdamore/tcell/v2"
)

func PresentFB(s tcell.Screen, fb *framebuffer.Framebuffer) {
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
	fmt.Printf("%d,%d\r\n", w, h)
	f1 := framebuffer.New(20, h-2)
	f1.Clear('.')

	f2 := framebuffer.New(20, h-2)
	f2.Clear(' ')
	p := Player{5, 5, '@'}
	f2.View[p.Y][p.X] = p.CH
	r := renderer.New(w, h)
	r.AddLayer(f1, 1, layer.BlendCopy)
	r.AddLayer(f2, 2, layer.BlendOr)

	defer screen.Fini()
	screen.Clear()
loop:
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEsc:
				break loop
			case tcell.KeyEnter:
				r.MarkDirty(layer.Rect{0, 0, w, h})
				r.Render()
				PresentFB(screen, r.OutputFront())
			}
		}
	}
}
