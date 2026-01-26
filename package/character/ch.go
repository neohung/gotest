// Package character
package character

type Avatar struct {
	x uint16
	y uint16
	c rune
}

func Init(x, y uint16, c rune) *Avatar {
	return &Avatar{
		x,
		y,
		c,
	}
}

func (a Avatar) Render(dst [][]rune) {
	dst[a.y][a.x] = a.c
}
