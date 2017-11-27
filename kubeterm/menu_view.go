package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"time"
)

type MenuView struct {
	gv    *gocui.View
	done  chan struct{}
	dirty bool
	items []MenuItem
}

type MenuItem interface {
	DisplayName() string
}

func NewMenuView() *MenuView {
	m := &MenuView{
		done:  make(chan struct{}),
		items: []MenuItem{},
		dirty: false,
	}

	return m
}

func (v *MenuView) Open(a *App, gv *gocui.View) {
	v.gv = gv
	a.SetViewKeybinding(v, 'q', ModNone, a.Quit)
	a.SetViewKeybinding(v, 'j', ModNone, v.ptrDown)
	a.SetViewKeybinding(v, 'k', ModNone, v.ptrUp)
	a.SetViewKeybinding(v, KeyEnter, ModNone, v.enter)
	gv.Highlight = true
	gv.SelBgColor = gocui.ColorRed
	gv.SelFgColor = gocui.ColorGreen

	tick := time.Tick(50 * time.Millisecond)

	v.draw()

	go func() {
		for {
			select {
			case <-v.done:
				return
			case <-tick:
				if v.dirty == true {
					v.dirty = false
					gv.Clear()
					v.draw()
				}
			default:
			}
		}
	}()
}

func (v *MenuView) enter() error {
	return nil
}

func (v *MenuView) AddMenu(item MenuItem) {
	v.items = append(v.items, item)
}

func (v *MenuView) Close() {
	close(v.done)
}

func (v *MenuView) draw() {
	for _, item := range v.items {
		fmt.Fprintln(v.gv, item.DisplayName())
	}
}

func (v *MenuView) Name() string {
	return "menu"
}

func (v *MenuView) ptrDown() error {
	x, y := v.gv.Cursor()
	next := y + 1

	if next > len(v.items)-1 {
		next = y
	}

	if err := v.gv.SetCursor(x, next); err != nil {
		return err
	}

	v.dirty = true

	return nil
}

func (v *MenuView) ptrUp() error {
	x, y := v.gv.Cursor()
	next := y - 1

	if next < 0 {
		next = 0
	}

	if err := v.gv.SetCursor(x, next); err != nil {
		return err
	}

	v.dirty = true

	return nil
}
