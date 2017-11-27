package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type MenuView struct {
	done  chan struct{}
	items []MenuItem
}

type MenuItem interface {
	DisplayName() string
	OnFocus(*App, *gocui.View)
	OnEnter(*App, *gocui.View)
	ViewContext
}

func NewMenuView() *MenuView {
	m := &MenuView{
		done:  make(chan struct{}),
		items: []MenuItem{},
	}

	return m
}

func (v *MenuView) Open(a *App, gv *gocui.View) {
	a.SetViewKeybinding(v, 'q', ModNone, a.Quit)
	a.SetViewKeybinding(v, 'j', ModNone, v.ptrDown)
	a.SetViewKeybinding(v, 'k', ModNone, v.ptrUp)
	a.SetViewKeybinding(v, KeyEnter, ModNone, v.enter)

	gv.SelBgColor = gocui.ColorRed
	gv.SelFgColor = gocui.ColorGreen

	v.Draw(gv, true)
}

func (v *MenuView) enter(a *App, gv *gocui.View) error {
	a.Update(func() error {
		v.Draw(gv, false)
		return nil
	})

	item := v.selectItem(a, gv)
	item.OnEnter(a, a.GetGoCuiView(item))
	a.SetCurrentView(item)

	return nil
}

func (v *MenuView) selectItem(a *App, gv *gocui.View) MenuItem {
	_, y := gv.Cursor()
	item := v.items[y]
	a.SetViewOnTop(item)
	item.OnFocus(a, a.GetGoCuiView(item))

	return item
}

func (v *MenuView) AddMenu(item MenuItem) {
	v.items = append(v.items, item)
}

func (v *MenuView) Close() {
	close(v.done)
}

func (v *MenuView) Draw(gv *gocui.View, hl bool) {
	gv.Highlight = hl
	gv.Clear()
	for _, item := range v.items {
		fmt.Fprintln(gv, item.DisplayName())
	}
}

func (v *MenuView) Name() string {
	return "menu"
}

func (v *MenuView) ptrDown(a *App, gv *gocui.View) error {
	x, y := gv.Cursor()
	next := y + 1

	if next > len(v.items)-1 {
		next = y
	}

	if err := gv.SetCursor(x, next); err != nil {
		return err
	}

	v.selectItem(a, gv)

	return nil
}

func (v *MenuView) ptrUp(a *App, gv *gocui.View) error {
	x, y := gv.Cursor()
	next := y - 1

	if next < 0 {
		next = 0
	}

	if err := gv.SetCursor(x, next); err != nil {
		return err
	}

	v.selectItem(a, gv)

	return nil
}
