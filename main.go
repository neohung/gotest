package main

import (
	layer "test/package/layer"
	maze "test/package/maze"
)

// 編譯期檢查
// var _ Renderable = (*MyData)(nil)
func main() {
	m := maze.New(20, 20)
	l := layer.New()
	l.Push(m)
	// l.Render()
}
