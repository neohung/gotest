// Package framebuffer
package framebuffer

// import "fmt"
type Framebuffer struct {
	W    int
	H    int
	Buf  []rune
	View [][]rune
}

func (fb *Framebuffer) Clear(ch rune) {
	for i := range fb.Buf {
		fb.Buf[i] = ch
	}
}

func New(w, h int) *Framebuffer {
	buf := make([]rune, w*h)
	view := make([][]rune, h)
	for y := range view {
		view[y] = buf[y*w : (y+1)*w]
	}
	return &Framebuffer{
		W:    w,
		H:    h,
		Buf:  buf,
		View: view,
	}
}
