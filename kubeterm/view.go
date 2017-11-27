package kubeterm

import (
	"github.com/jroimartin/gocui"
	"log"
)

type View interface {
	Name() string
	Open(a *App, v *gocui.View)
	Close()
}

func (a *App) AppendView(v View) {
	a.views = append(a.views, v)
}

func (a *App) OpenView(v View, x0, y0, x1, y1 int) {
	view, err := a.g.SetView(v.Name(), x0, y0, x1, y1)

	view.Frame = false

	if err != nil &&
		err != gocui.ErrUnknownView {
		log.Panicln(err)
	}

	v.Open(a, view)
	a.AppendView(v)
}

func (a *App) SetCurrentView(v View) {
	_, err := a.g.SetCurrentView(v.Name())

	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
}

func (a *App) SetViewKeybinding(v View, key interface{}, mod gocui.Modifier, handler func(*App, *gocui.View) error) {
	a.setKeybinding(v.Name(), key, mod, handler)
}

func (a *App) setKeybinding(viewname string, key interface{}, mod gocui.Modifier, handler func(*App, *gocui.View) error) {
	f := func(g *gocui.Gui, v *gocui.View) error {
		return handler(a, v)
	}

	if err := a.g.SetKeybinding(viewname, key, mod, f); err != nil {
		log.Panicln(err)
	}
}

func (a *App) SetViewOnTop(v View) {
	_, err := a.g.SetViewOnTop(v.Name())

	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
}

func (a *App) Update(handler func() error) {
	f := func(*gocui.Gui) error { return handler() }
	a.g.Update(f)
}

func (a *App) GetGoCuiView(v View) *gocui.View {
	gv, err := a.g.View(v.Name())

	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}

	return gv
}
