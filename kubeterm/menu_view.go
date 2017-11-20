package kubeterm

import (
	"github.com/jroimartin/gocui"
)

type Drawer interface {
	Draw(view *gocui.View)
}

type MenuView struct {
	drawer Drawer
	client *Client
}

func NewMenuView(c *Client) *MenuView {
	return &MenuView{
		drawer: initialDrawer(c),
		client: c,
	}
}

func initialDrawer(c *Client) Drawer {
	return NewNamespaceMode(c)
}

func (m *MenuView) Draw(v *View) {
	vctx := v.SetView("namespace", -1, -1, v.width/4, v.height/3)

	m.drawer.Draw(vctx)
}
