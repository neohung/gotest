// Package layer
package layer

import (
	"container/list"
)

type Renderable interface {
	Render(dst [][]rune)
}

type Layer struct {
	// Renderable
	ll *list.List
}

func New() *Layer {
	l := Layer{
		ll: list.New(),
	}
	return &l
}

func (l *Layer) Push(e Renderable) {
	l.ll.PushBack(e)
}

func (l *Layer) Render(dst [][]rune) {
	// 型別斷言
	for e := l.ll.Front(); e != nil; e = e.Next() {
		if v, ok := e.Value.(Renderable); ok {
			v.Render(dst)
		}
	}
}
