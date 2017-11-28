package kubeterm

import (
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	"log"
	"sync"
)

type Window struct {
	gui *gocui.Gui

	Width       int
	Height      int
	Views       []*View
	CurrentView *View
	stack       *stack
}

var instance *Window = newWindow()

func GetWindow() *Window {
	return instance
}

type stack struct {
	lock sync.Mutex
	s    []*View
}

func newStack() *stack {
	return &stack{sync.Mutex{}, make([]*View, 0)}
}

func (s *stack) Pop() (*View, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return nil, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) Push(v *View) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) IsEmpty() bool {
	return len(s.s) == 0
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
		stack:  newStack(),
	}
}

func (w *Window) Back() error {
	if w.stack.IsEmpty() {
		return w.Quit()
	}

	v, err := w.stack.Pop()

	if err != nil {
		return err
	}

	w.SetCurrentView(v)

	return nil
}

func (w *Window) Quit() error {
	for _, view := range w.Views {
		if err := view.Quit(); err != nil {
			log.Panicln(err)
		}
	}
	return gocui.ErrQuit
}

func (w *Window) SetViewOnTop(v *View) {
	_, err := w.gui.SetViewOnTop(v.Name())

	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
}

func (w *Window) SetCurrentView(v *View) {
	if w.CurrentView != nil {
		w.CurrentView.OnFocusOut()
		w.stack.Push(w.CurrentView)
	}

	_, err := w.gui.SetCurrentView(v.Name())

	v.OnFocusIn()
	w.CurrentView = v

	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
}

func (w *Window) AddView(ctx ViewContext) *View {
	x0, y0, x1, y1 := ctx.Position()
	v, err := w.gui.SetView(ctx.Name(), x0, y0, x1, y1)
	if err != nil &&
		err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
	v.Frame = false
	view := NewView(w.gui, ctx, v)

	w.Views = append(w.Views, view)

	return view
}

func (w *Window) Init() {
	err := w.gui.SetKeybinding("", KeyCtrlC, gocui.ModNone, func(*gocui.Gui, *gocui.View) error {
		return w.Quit()
	})

	for _, view := range w.Views {
		view.Init()
	}

	if err != nil {
		log.Panicln(err)
	}
}

func (w *Window) Loop() {
	defer w.gui.Close()

	w.Init()

	for _, view := range w.Views {
		view.Update()
	}

	if err := w.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
