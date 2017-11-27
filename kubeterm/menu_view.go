package kubeterm

type MenuView struct {
	items []*View
}

func NewMenuView() *MenuView {
	m := &MenuView{
		items: []*View{},
	}

	return m
}

func (v *MenuView) DisplayName() string {
	return "Menu"
}

func (v *MenuView) Init(view *View) {
	view.SetKeybinding('q', func() error {
		return GetWindow().Quit()
	})
	view.SetKeybinding('j', func() error {
		if err := view.PointerDown(); err != nil {
			return err
		}

		ptr := view.PointerPos()
		item := v.items[ptr]
		GetWindow().SetViewOnTop(item)
		return nil
	})

	view.SetKeybinding('k', func() error {
		if err := view.PointerUp(); err != nil {
			return err
		}

		ptr := view.PointerPos()
		item := v.items[ptr]
		GetWindow().SetViewOnTop(item)
		return nil
	})

	view.SetKeybinding(KeyEnter, func() error {
		ptr := view.PointerPos()
		item := v.items[ptr]
		GetWindow().SetCurrentView(item)
		return nil
	})
}

func (v *MenuView) BeginPointerIndex() (x int) {
	return 0
}

func (v *MenuView) Position() (x0, y0, x1, y1 int) {
	return 0, 0, 20, GetWindow().Height
}

func (v *MenuView) Height() int {
	return len(v.items)
}

func (v *MenuView) Lines() []string {
	lines := []string{}
	for _, item := range v.items {
		lines = append(lines, item.DisplayName())
	}
	return lines
}

func (v *MenuView) AddMenu(item *View) {
	v.items = append(v.items, item)
}

func (v *MenuView) Name() string {
	return "menu"
}
