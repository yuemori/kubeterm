package kubeterm

import ()

type StatusView struct {
	drawer Drawer
	client *Client
}

func NewStatusView(c *Client) *StatusView {
	return &StatusView{
		drawer: initialStatusDrawer(c),
		client: c,
	}
}

func initialStatusDrawer(c *Client) Drawer {
	return NewPodMode(c, "kube-system")
}

func (m *StatusView) Draw(view *View) {
	v := view.SetView("status", view.width/4, -1, view.width, view.height/3)

	m.drawer.Draw(v)
}
