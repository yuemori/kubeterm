package kubeterm

import (
	"github.com/jroimartin/gocui"
	"log"
)

type Window struct {
	gui *gocui.Gui

	Width  int
	Height int
}

func NewWindow() *Window {
	g, err := gocui.NewGui(gocui.OutputNormal)
	w, h := g.Size()

	if err != nil {
		log.Panicln(err)
	}

	return &Window{
		gui:    g,
		Width:  w,
		Height: h,
	}
}

// func (w *Window) Loop() {
// 	defer w.gui.Close()
//
// 	if err := w.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
// 		log.Panicln(err)
// 	}
// }
//
// func (w *Window) SetCurrentView(v ViewContext) {
// 	_, err := w.gui.SetCurrentView(v.Name())
//
// 	if err != nil && err != gocui.ErrUnknownView {
// 		log.Panicln(err)
// 	}
// }
