// Package layer
package layer

import (
	"container/list"
	"fmt"
)

type Renderable interface {
	Render() string
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

func (l *Layer) Render() {
	// 型別斷言
	for e := l.ll.Front(); e != nil; e = e.Next() {
		if v, ok := e.Value.(Renderable); ok {
			fmt.Printf("%s", v.Render())
		}
	}
}
