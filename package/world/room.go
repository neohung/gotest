package world

import "math/rand"

type Tile rune

const (
	Wall       Tile = '#'
	Floor      Tile = '.'
	Grass      Tile = '"'
	Water      Tile = '~'
	SpawnPoint Tile = 'S'
)

type Room struct {
	X, Y, W, H int
}

func (r Room) Center() (int, int) {
	return r.X + r.W/2, r.Y + r.H/2
}

func CarveRoom(w *World, r Room) {
	for y := r.Y; y < r.Y+r.H; y++ {
		for x := r.X; x < r.X+r.W; x++ {
			w.Tiles[y][x] = rune(Floor)
		}
	}
}

func carveHCarridor(w *World, x1, x2, y int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for x := x1; x <= x2; x++ {
		w.Tiles[y][x] = rune(Floor)
	}
}

func carveVCarridor(w *World, y1, y2, x int) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for y := y1; y <= y2; y++ {
		w.Tiles[y][x] = rune(Floor)
	}
}

// CarveCorridor L型走廊
func CarveCorridor(w *World, x1, y1, x2, y2 int) {
	if rand.Intn(2) == 0 {
		carveHCarridor(w, x1, x2, y1)
		carveVCarridor(w, y1, y2, x2)
	} else {
		carveVCarridor(w, y1, y2, x1)
		carveHCarridor(w, x1, x2, y2)
	}
}

// SprinkleNature 添加自然地形, Floor 5%機率變成Gress, 3%機率變Water
func SprinkleNature(w *World) {
	for y := 1; y < w.H-1; y++ {
		for x := 1; x < w.W-1; x++ {
			if w.Tiles[y][x] == rune(Floor) {
				glassORwater := rand.Float64()
				if glassORwater < 0.05 {
					w.Tiles[y][x] = rune(Grass)
				} else if glassORwater < 0.08 {
					w.Tiles[y][x] = rune(Water)
				}
			}
		}
	}
}
