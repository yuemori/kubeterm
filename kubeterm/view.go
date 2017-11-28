package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

type View struct {
	gui  *gocui.Gui
	ctx  ViewContext
	view *gocui.View

	keepPointer bool
}

type ViewContext interface {
	Name() string
	BeginPointerIndex() int
	Init(view *View)
	Position() (int, int, int, int)
	Lines() []string
	Height() int
	OnNamespaceUpdate(string)
	OnVisible()
	OnInvisible()
	OnFocusOut()
	OnFocusIn()
}

func NewView(gui *gocui.Gui, ctx ViewContext, view *gocui.View) *View {
	return &View{gui, ctx, view, false}
}

func (v *View) Name() string {
	return v.ctx.Name()
}

func (v *View) KeepPointerOnFocusOut() {
	v.keepPointer = true
}

func (v *View) SetKeybinding(key interface{}, handler func() error) {
	f := func(g *gocui.Gui, v *gocui.View) error {
		return handler()
	}

	if err := v.gui.SetKeybinding(v.ctx.Name(), key, gocui.ModNone, f); err != nil {
		log.Panicln(err)
	}
}

func (v *View) Init() {
	v.view.SelBgColor = gocui.ColorRed
	v.view.SelFgColor = gocui.ColorGreen
	v.ctx.Init(v)
}

func (v *View) PointerUp() error {
	x, y := v.view.Cursor()
	next := y - 1

	if next < v.ctx.BeginPointerIndex() {
		next = v.ctx.BeginPointerIndex()
	}

	if err := v.view.SetCursor(x, next); err != nil {
		return err
	}

	return nil
}

func (v *View) PointerDown() error {
	x, y := v.view.Cursor()
	next := y + 1

	if next > v.ctx.Height()-1 {
		next = v.ctx.Height() - 1
	}

	if err := v.view.SetCursor(x, next); err != nil {
		return err
	}

	return nil
}

func (v *View) PointerPos() int {
	_, y := v.view.Cursor()

	return y
}

func (v *View) pointerReset() {
	if err := v.view.SetCursor(0, v.ctx.BeginPointerIndex()); err != nil {
		log.Panicln(err)
	}
}

func (v *View) OnVisible() {
	v.ctx.OnVisible()
	v.Update()
}

func (v *View) OnNamespaceUpdate(namespace string) {
	v.ctx.OnNamespaceUpdate(namespace)
}

func (v *View) OnInvisible() {
	v.ctx.OnInvisible()
	v.Update()
}

func (v *View) OnFocusIn() {
	v.ctx.OnFocusIn()
	if !v.keepPointer {
		v.pointerReset()
	}
	v.view.Highlight = true
	v.Update()
}

func (v *View) OnFocusOut() {
	v.ctx.OnFocusOut()
	v.view.Highlight = false
	v.Update()
}

func (v *View) Quit() error {
	return GetWindow().Back()
}

func (v *View) Cursor() (x, y int) {
	return v.view.Cursor()
}

func (v *View) Update() {
	v.gui.Update(func(*gocui.Gui) error {
		v.view.Clear()
		for _, line := range v.ctx.Lines() {
			fmt.Fprintln(v.view, line)
		}
		return nil
	})
}
