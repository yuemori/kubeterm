package kubeterm

import (
	"github.com/nsf/termbox-go"
	"time"
)

type Mode interface {
	Draw() error
}

type View struct {
	quit   bool
	width  int
	height int
	top    int64
	ptr    int64
	mode   Mode
}

func NewView(client *Client) *View {
	w, h := termbox.Size()
	mode := &NamespaceMode{}
	mode.Init(client)

	return &View{
		width:  w,
		height: h,
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
		}
	}
}

func (v *View) draw() {
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				termbox.HideCursor()
				if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
					panic(err)
				}

				if err := v.mode.Draw(); err != nil {
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
