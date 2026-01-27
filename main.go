package main

import (
	//"fmt"
	"os"
	"time"
	//"strings"

	// tea "github.com/charmbracelet/bubbletea"
	//"fmt"

	"test/package/roguelike/framebuffer"
	"test/package/roguelike/layer"
	"test/package/roguelike/renderer"

	"github.com/gdamore/tcell/v2"
)

type Player struct {
	X, Y int
	CH   rune
}

func PresentFB(screen tcell.Screen, fb *framebuffer.Framebuffer) {
	for y := 0; y < fb.H; y++ {
		for x := 0; x < fb.W; x++ {
			ch := fb.View[y][x]
			screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}
	screen.Show()
}

func genMap(fb *framebuffer.Framebuffer) {
	fb.Clear('.')
	for x := range fb.W {
		fb.View[0][x] = '#'
		fb.View[fb.H-1][x] = '#'
	}
	for y := range fb.H {
		fb.View[y][0] = '#'
		fb.View[y][fb.W-1] = '#'
	}
}

func main() {
	// 建立screen
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	w, h := screen.Size()
	r := renderer.NewRenderer(w, h)

	mapFB := framebuffer.New(w, h)
	actorFB := framebuffer.New(w, h)
	hudFB := framebuffer.New(w, h)

	r.AddLayer(mapFB, 0, layer.BlendCopy)
	r.AddLayer(actorFB, 1, layer.BlendOr)
	r.AddLayer(hudFB, 2, layer.BlendOr)

	genMap(mapFB)

	// 初始化後
	// 將 mapFB 先渲染到 front/back buffer，保證第一 frame 就完整
	//copyRect(r.front, mapFB, Rect{0, 0, w, h})
	//`pyRect(r.back, mapFB, Rect{0, 0, w, h})

	player := Player{X: w / 2, Y: h / 2, CH: '@'}

	// 初始 dirty
	r.MarkDirty(layer.Rect{
		X: 0,
		Y: 0,
		W: w,
		H: h,
	})

	screen.EnableMouse()
loop:
	for {
		screen.Clear()
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// 存下player舊座標等等make dirty使用
			oldX, oldY := player.X, player.Y
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				break loop
				// return
			case tcell.KeyUp:
				player.Y--
			case tcell.KeyDown:
				player.Y++
			case tcell.KeyLeft:
				player.X--
			case tcell.KeyRight:
				player.X++
			}

			r.MarkDirty(layer.Rect{X: oldX, Y: oldY, W: 1, H: 1})
			r.MarkDirty(layer.Rect{X: player.X, Y: player.Y, W: 1, H: 1})
		case *tcell.EventResize:
			w, h = screen.Size()
			// screen.Sync()
		}
		//================================================
		// Actor
		actorFB.Clear(' ')
		actorFB.View[player.Y][player.X] = player.CH
		//================================================
		// HUD
		hudFB.Clear(' ')
		msg := "Move: Arrow / WASD | ESC quit"
		for i, ch := range msg {
			hudFB.View[h-1][i] = ch
		}
		r.MarkDirty(layer.Rect{
			X: 0,
			Y: h - 1,
			W: w,
			H: 1,
		})
		//================================================
		r.Render()
		// render完後, 取出front framebuffer, present到screen
		PresentFB(screen, r.Front())
		time.Sleep(16 * time.Millisecond)
	}
	screen.Fini()
	os.Exit(0)
}
