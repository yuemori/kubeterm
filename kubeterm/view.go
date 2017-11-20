package kubeterm

import (
	"github.com/jroimartin/gocui"
	"log"
)

type View struct {
	width  int
	height int
	top    int
	ptr    int
	g      *gocui.Gui
}

func NewView(client *Client) *View {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	w, h := g.Size()

	return &View{
		width:  w,
		height: h,
		top:    0,
		ptr:    0,
		g:      g,
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (v *View) ptrDown(_ *gocui.Gui, _ *gocui.View) error {
	v.ptr++
	v.fixPtr()

	return nil
}

func (v *View) ptrUp(_ *gocui.Gui, _ *gocui.View) error {
	v.ptr--
	v.fixPtr()

	return nil
}

func (v *View) Loop(client *Client) {
	defer v.g.Close()

	v.registerKeyBindings()

	menu := NewMenuView(client)
	menu.Draw(v)

	status := NewStatusView(client)
	status.Draw(v)

	if err := v.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (v *View) SetKeybinding(viewname string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
	if err := v.g.SetKeybinding(viewname, key, mod, handler); err != nil {
		log.Panicln(err)
	}
}

func (v *View) SetView(name string, x0, y0, x1, y1 int) *gocui.View {
	view, err := v.g.SetView(name, x0, y0, x1, y1)

	if err != nil &&
		err != gocui.ErrUnknownView {
		log.Panicln(err)
	}

	return view
}

func (v *View) registerKeyBindings() {
	v.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	v.SetKeybinding("", 'q', gocui.ModNone, quit)
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
