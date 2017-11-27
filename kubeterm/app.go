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
	Client *Client
	menu   *MenuView

	Window           *Window
	currentNamespace string

	views []ViewContext
}

func NewApp(client *Client) *App {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	app := &App{
		Client: client,
		Window: GetWindow(),
		g:      g,
		views:  []ViewContext{},
	}

	return app
}

func (app *App) MainLoop() {
	ns := NewNamespaceView(app.Client)
	pod := NewPodView("default", app.Client)

	menu := NewMenuView()
	menu.AddMenu(ns)
	menu.AddMenu(pod)

	app.Window.AddView(ns)
	app.Window.AddView(pod)
	app.Window.AddView(menu)

	app.Window.DisplayView(ns)

	app.Window.Loop()
}

func (app *App) Quit() error {
	return app.Window.Quit()
}

func (a *App) Update(handler func() error) {
	f := func(*gocui.Gui) error { return handler() }
	a.g.Update(f)
}

func (a *App) GetGoCuiView(v ViewContext) *gocui.View {
	gv, err := a.g.View(v.Name())

	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}

	return gv
}
