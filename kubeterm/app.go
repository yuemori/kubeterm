package kubeterm

import (
	"github.com/jroimartin/gocui"
)

const (
	KeyCtrlC = gocui.KeyCtrlC
	KeyEnter = gocui.KeyEnter
)

const (
	ModNone = gocui.ModNone
	ModAlt  = gocui.ModAlt
)

const (
	DEFAULT_NAMESPACE = "default"
)

var app *App = NewApp()

type App struct {
	CurrentNamespace string
}

func GetApp() *App {
	return app
}

func NewApp() *App {
	app := &App{DEFAULT_NAMESPACE}

	return app
}

func (app *App) MainLoop(client *Client) {
	w := GetWindow()
	ns := w.AddView(NewNamespaceView(client))
	pod := w.AddView(NewPodView(client))
	menu := NewMenuView()
	menu.AddMenu(ns)
	menu.AddMenu(pod)

	w.SetViewOnTop(ns)
	w.SetCurrentView(w.AddView(menu))

	w.Loop()
}
func (app *App) SetCurrentNamespace(namespace string) {
	app.CurrentNamespace = namespace
	GetWindow().OnNamespaceUpdate(namespace)
}
