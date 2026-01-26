package main

import (
	layer "test/package"
)

type MyData struct {
	s string
}

func (m MyData) Render() string {
	return m.s
}

// 編譯期檢查
// var _ Renderable = (*MyData)(nil)
func main() {
	a1 := MyData{s: "A"}
	a2 := MyData{s: "B"}
	a3 := MyData{s: "C"}
	l := layer.New()
	l.Push(a1)
	l.Push(a2)
	l.Push(a3)
	l.Render()
}
