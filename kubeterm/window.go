package kubeterm

import (
	"github.com/jroimartin/gocui"
	"log"
)

type Window struct {
	gui *gocui.Gui

	Width  int
	Height int
	Views  []*View
}

var instance *Window = newWindow()

func GetWindow() *Window {
	return instance
}

func newWindow() *Window {
	g, err := gocui.NewGui(gocui.OutputNormal)
	w, h := g.Size()

	if err != nil {
		log.Panicln(err)
	}

	return &Window{
		gui:    g,
		Width:  w,
		Height: h,
		Views:  []*View{},
	}
}

func (w *Window) Quit() error {
	for _, view := range w.Views {
		if err := view.Quit(); err != nil {
			log.Panicln(err)
		}
	}
	return nil
}

func (w *Window) DisplayView(ctx ViewContext) error {
	_, err := w.gui.SetViewOnTop(ctx.Name())
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func (w *Window) FocusIn(ctx ViewContext) error {
	_, err := w.gui.SetCurrentView(ctx.Name())

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func (w *Window) AddView(ctx ViewContext) {
	x0, y0, x1, y1 := ctx.Position()
	v, err := w.gui.SetView(ctx.Name(), x0, y0, x1, y1)
	if err != nil &&
		err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
	v.Frame = false
	view := NewView(w.gui, ctx, v)

	w.Views = append(w.Views, view)
}

func (w *Window) Loop() {
	defer w.gui.Close()

	w.gui.SetKeybinding("", KeyCtrlC, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		w.Quit()
		return nil
	})

	if err := w.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
