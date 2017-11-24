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
	g      *gocui.Gui
	client *Client

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
		client:    client,
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
	a.openMenuView()
	a.openNamespaceView()

	if err := a.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) openMenuView() {
	v := NewMenuView()
	a.OpenView(v, 0, 0, 20, a.MaxHeight)
	a.SetCurrentView(v)
}

func (a *App) openNamespaceView() {
	v := NewNamespaceView()
	a.OpenView(v, 20, 0, a.MaxWidth, a.MaxHeight)
}

func (a *App) Quit() error {
	for _, v := range a.views {
		v.Close()
	}

	return gocui.ErrQuit
}