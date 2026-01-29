// Package world
package world

import (
	"math/rand"

	"test2/package/framebuffer"
)

type World struct {
	W, H  int
	Tiles [][]rune
	Cam   Camera
}

func New(w, h int, c Camera) *World {
	world := make([]rune, w*h)
	tiles := make([][]rune, h)
	for y := range h {
		tiles[y] = world[y*w : (y+1)*w]
		for x := range w {
			// tiles[y][x] = '.'
			tiles[y][x] = rune(Wall)
		}
	}
	return &World{
		W:     w,
		H:     h,
		Tiles: tiles,
		Cam:   c,
	}
}

func (w *World) FindSpawn() (x, y int) {
	for y = 0; y < w.H; y++ {
		for x = 0; x < w.W; x++ {
			if w.Tiles[y][x] == rune(SpawnPoint) {
				w.Tiles[y][x] = rune(Floor) // 清掉標記
				return
			}
		}
	}
	panic("no spawn point")
}

func (w *World) GenerateMap() {
	var rooms []Room
	maxRooms := 8
	for range maxRooms {
		rw := rand.Intn(6) + 4
		rh := rand.Intn(6) + 4
		rx := rand.Intn(w.W-rw-2) + 1
		ry := rand.Intn(w.H-rh-2) + 1

		room := Room{rx, ry, rw, rh}
		// 避免重疊
		ok := true
		for _, other := range rooms {
			if rx < other.X+other.W &&
				rx+rw > other.X &&
				ry < other.Y+other.H &&
				ry+rh > other.Y {
				ok = false
				break
			}
		}
		if ok {
			CarveRoom(w, room)
			if len(rooms) > 0 {
				// 連接兩個房間的通道
				x1, y1 := rooms[len(rooms)-1].Center()
				x2, y2 := room.Center()
				CarveCorridor(w, x1, y1, x2, y2)
			}
			rooms = append(rooms, room)
		}
	}
	SprinkleNature(w)
	// Spawn In rooms
	sx, sy := rooms[0].Center()
	w.Tiles[sy][sx] = rune(SpawnPoint)
}

func (w *World) RenderToMapFB(fb *framebuffer.Framebuffer) {
	for sy := range w.Cam.H {
		// 遍歷camera區域在world上的y
		wy := w.Cam.Y + sy
		if wy < 0 || wy >= w.H {
			continue
		}
		for sx := range w.Cam.W {
			wx := w.Cam.X + sx
			if wx < 0 || wx >= w.W {
				continue
			}
			// Render world tile到mapFB
			fb.View[sy][sx] = w.Tiles[wy][wx]
		}
	}
}
