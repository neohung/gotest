// Package world
package world

import (
	"math/rand"

	"test2/package/framebuffer"
)

type World struct {
	W, H    int
	Tiles   [][]rune
	Cam     Camera
	Visible [][]bool // 本回合可見
	Seen    [][]bool // 曾經看過（fog of war）
}

func New(w, h int, c Camera) *World {
	world := make([]rune, w*h)
	tiles := make([][]rune, h)
	vis := make([][]bool, h)
	seen := make([][]bool, h)
	for y := range h {
		tiles[y] = world[y*w : (y+1)*w]
		vis[y] = make([]bool, w)
		seen[y] = make([]bool, w)
		for x := range w {
			// tiles[y][x] = '.'
			tiles[y][x] = rune(Wall)
		}
	}
	return &World{
		W:       w,
		H:       h,
		Tiles:   tiles,
		Cam:     c,
		Visible: vis,
		Seen:    seen,
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
			ch := ' ' // 看不到顯示空格
			// Render world tile到mapFB
			if w.Visible[wy][wx] {
				ch = w.Tiles[wy][wx]
			} else if w.Seen[wy][wx] {
				ch = rune(Floor)
			}
			fb.View[sy][sx] = ch
		}
	}
}

func (w *World) BlocksSight(x, y int) bool {
	if x < 0 || y < 0 || x >= w.W || y >= w.H {
		return true
	}
	return w.Tiles[y][x] == rune(Wall)
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Bresenham(x0, y0, x1, y1 int, fn func(x, y int) bool) {
	dx := absInt(x1 - x0)
	dy := -absInt(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy

	for {
		if !fn(x0, y0) {
			return
		}
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func (w *World) ComputeFOV(px, py int, radius int) {
	// 清空可見 重新計算
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			w.Visible[y][x] = false
		}
	}
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			tx := px + dx
			ty := py + dy

			if dx*dx+dy*dy > radius*radius {
				continue
			}

			Bresenham(px, py, tx, ty, func(x, y int) bool {
				if x < 0 || y < 0 || x >= w.W || y >= w.H {
					return false
				}
				w.Visible[y][x] = true
				w.Seen[y][x] = true
				return !w.BlocksSight(x, y)
			})
		}
	}
}
