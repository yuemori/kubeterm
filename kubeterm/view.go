package kubeterm

import (
	"github.com/nsf/termbox-go"
	"time"
)

type Mode interface {
	Draw(v *View) error
}

type View struct {
	quit   bool
	width  int
	height int
	top    int
	ptr    int
	mode   Mode
}

func NewView(client *Client) *View {
	w, h := termbox.Size()
	mode := &NamespaceMode{}
	mode.Init(client)

	return &View{
		width:  w,
		height: h,
		top:    0,
		ptr:    0,
		mode:   mode,
	}
}

func (v *View) Loop() {
	evCh := make(chan termbox.Event)
	go func() {
		for {
			evCh <- termbox.PollEvent()
		}
	}()

	v.draw()

	tick := time.Tick(time.Second / 2)
	for {
		select {
		case <-tick:
		case ev := <-evCh:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyCtrlC {
				return
			}
			v.updateEvent(ev)
		}

		if v.quit {
			break
		}
	}
}

func (v *View) updateEvent(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventResize:
		v.width, v.height = ev.Width, ev.Height
		// v.fixPtr()
	case termbox.EventKey:
		switch ev.Ch {
		case 'q':
			v.quit = true
		case 'j':
			v.ptr++
			v.fixPtr()
		case 'k':
			v.ptr--
			v.fixPtr()
		}
	}
}

func (v *View) fixPtr() {
	if v.ptr < 0 {
		v.ptr = 0
	}

	if v.ptr < v.top {
		v.top = v.ptr
		return
	}

	end := v.calcEnd()
	if v.ptr >= end {
		v.top += v.ptr - end + 1
		return
	}
}

func (v *View) calcEnd() int {
	h := v.height - 2
	if h < 1 {
		h = 1
	}

	end := v.top + h
	return end
}

func (v *View) draw() {
	ticker := time.NewTicker(50 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				termbox.HideCursor()
				if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
					panic(err)
				}

				if err := v.mode.Draw(v); err != nil {
					panic(err)
				}

				if err := termbox.Flush(); err != nil {
					panic(err)
				}

			default:
				// nothing to do
			}
		}
	}()
}
