package kubeterm

import (
	"github.com/jroimartin/gocui"
	"log"
)

type View struct {
	Width  int
	Height int
	g      *gocui.Gui
}

func NewView(client *Client) *View {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	w, h := g.Size()

	return &View{
		Width:  w,
		Height: h,
		g:      g,
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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

	view.Frame = false

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
