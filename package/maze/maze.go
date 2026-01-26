// Package maze
package maze

import "unsafe"

type Maze struct {
	context [][]rune
}

func New(width, height uint16) *Maze {
	m := Maze{
		context: make([][]rune, height),
	}
	tmp := make([]rune, width*height)
	for i := range m.context {
		m.context[i] = tmp[uint16(i)*width : uint16(i+1)*width]
	}
	return &m
}

func (m *Maze) Maze1() {
	for i := range len(m.context) {
		for j := range len(m.context[0]) {
			if i == 0 || i == 19 || j == 0 || j == 19 {
				m.context[i][j] = '#'
			} else {
				m.context[i][j] = ' '
			}
		}
	}
}

func (m *Maze) Render(dst [][]rune) {
	dstView := unsafe.Slice(&dst[0][0], len(dst)*len(dst[0]))
	srcView := unsafe.Slice(&m.context[0][0], len(m.context)*len(m.context[0]))
	copy(dstView, srcView)
}
