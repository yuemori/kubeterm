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
		Window: NewWindow(),
		g:      g,
		views:  []ViewContext{},
	}

	return app
}

func (a *App) MainLoop() {
	defer a.g.Close()

	var items []MenuItem

	a.SetCurrentNamespace("default")
	a.setKeybinding("", KeyCtrlC, ModNone, a.Quit)
	nv := a.openNamespaceView()
	items = append(items, nv)
	items = append(items, a.openPodView())
	a.SetViewOnTop(nv)
	a.menu = a.openMenuView(items)

	if err := a.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (a *App) ReturnToMenu() {
	a.SetCurrentView(a.menu)
	a.menu.Draw(a.GetGoCuiView(a.menu), true)
}

func (a *App) CurrentNamespace() string {
	return a.currentNamespace
}

func (a *App) SetCurrentNamespace(ns string) {
	a.currentNamespace = ns
}

func (a *App) openMenuView(items []MenuItem) *MenuView {
	v := NewMenuView()
	for _, item := range items {
		v.AddMenu(item)
	}

	a.OpenView(v, 0, 0, 20, a.Window.Height)
	a.SetCurrentView(v)

	return v
}

func (a *App) openNamespaceView() *NamespaceView {
	v := NewNamespaceView(a.Client)
	a.OpenView(v, 20, 0, a.Window.Width, a.Window.Height)
	return v
}

func (a *App) openPodView() *PodView {
	v := NewPodView(a.CurrentNamespace(), a.Client)
	a.OpenView(v, 20, 0, a.Window.Width, a.Window.Height)
	return v
}

func (a *App) Quit(*App, *gocui.View) error {
	for _, v := range a.views {
		v.Close()
	}

	return gocui.ErrQuit
}
