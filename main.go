package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	avatar "test/package/character"
	layer "test/package/layer"
	maze "test/package/maze"
)

type MyData struct {
	context [][]rune
	l       *layer.Layer
}

func (m MyData) Init() tea.Cmd {
	return nil
}

func (m MyData) Dump() string {
	var sb strings.Builder
	for i := range len(m.context) {
		sb.WriteString(string(m.context[i][:]))
		sb.WriteByte('\n')
	}
	sb.WriteString("Press q to Quit")
	return sb.String()
}

func (m MyData) View() string {
	m.l.Render(m.context)
	return m.Dump()
}

func (m MyData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgtype := msg.(type) {
	case tea.KeyMsg:
		switch msgtype.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func initModel() *MyData {
	width := 20
	height := 20
	m := MyData{
		l:       layer.New(),
		context: make([][]rune, height),
	}
	tmp := make([]rune, width*height)
	for i := range m.context {
		m.context[i] = tmp[i*width : (i+1)*width]
	}
	mymaze := maze.New(uint16(width), uint16(height))
	mymaze.Maze1()
	m.l.Push(mymaze)

	a := avatar.Init(5, 5, '@')
	m.l.Push(a)
	return &m
}

// 編譯期檢查
// var _ Renderable = (*MyData)(nil)
func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("error occur")
		os.Exit(1)
	}
}
