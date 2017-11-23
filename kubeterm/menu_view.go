package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"time"
)

var (
	MenuLines = []string{
		"Namespaces",
		"Pods",
		"Services",
		"Deployments",
		"StatefulSets",
	}
)

type MenuView struct {
	v     *gocui.View
	done  chan struct{}
	dirty bool
}

func NewMenuView() *MenuView {
	m := &MenuView{
		done:  make(chan struct{}),
		dirty: false,
	}

	return m
}

func (m *MenuView) Open(a *App, v *gocui.View) {
	m.v = v
	a.SetViewKeybinding(m, 'q', ModNone, a.Quit)
	a.SetViewKeybinding(m, 'j', ModNone, m.ptrDown)
	a.SetViewKeybinding(m, 'k', ModNone, m.ptrUp)
	a.SetViewKeybinding(m, KeyEnter, ModNone, m.enter)
	v.Highlight = true
	v.SelBgColor = gocui.ColorRed
	v.SelFgColor = gocui.ColorGreen

	tick := time.Tick(50 * time.Millisecond)

	m.draw()

	go func() {
		for {
			select {
			case <-m.done:
				return
			case <-tick:
				if m.dirty == true {
					m.dirty = false
					m.clear()
					m.draw()
				}
			default:
			}
		}
	}()
}

func (m *MenuView) enter() error {
	return nil
}

func (m *MenuView) Close() {
	close(m.done)
}

func (m *MenuView) clear() {
	m.v.Clear()
}

func (m *MenuView) draw() {
	for _, str := range MenuLines {
		fmt.Fprintln(m.v, str)
	}
}

func (m *MenuView) Name() string {
	return "menu"
}

func (m *MenuView) ptrDown() error {
	x, y := m.v.Cursor()
	next := y + 1

	if next > len(MenuLines)-1 {
		next = y
	}

	if err := m.v.SetCursor(x, next); err != nil {
		return err
	}

	m.dirty = true

	return nil
}

func (m *MenuView) ptrUp() error {
	x, y := m.v.Cursor()
	next := y - 1

	if next < 0 {
		next = 0
	}

	if err := m.v.SetCursor(x, next); err != nil {
		return err
	}

	m.dirty = true

	return nil
}
