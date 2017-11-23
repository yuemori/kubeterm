package kubeterm

import (
	"github.com/jroimartin/gocui"
	"log"
)

const (
	KeyCtrlC = gocui.KeyCtrlC
	KeyEnter = gocui.KeyEnter
)

const (
	ModNone = gocui.ModNone
	ModAlt  = gocui.ModAlt
)

type App struct {
	g *gocui.Gui

	MaxHeight int
	MaxWidth  int

	views []View
}

func NewApp(client *Client) *App {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	w, h := g.Size()
	app := &App{
		MaxHeight: h,
		MaxWidth:  w,
		g:         g,
		views:     []View{},
	}

	return app
}

func (a *App) MainLoop() {
	defer a.g.Close()

	a.setKeybinding("", KeyCtrlC, ModNone, a.Quit)
	a.openMenu()

	if err := a.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) openMenu() {
	m := NewMenuView()
	a.OpenView(m, 0, 0, a.MaxWidth/3, a.MaxHeight)
	a.SetCurrentView(m)
}

func (a *App) Quit() error {
	for _, v := range a.views {
		v.Close()
	}

	return gocui.ErrQuit
}
