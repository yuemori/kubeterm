package kubeterm

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type NamespaceMode struct {
	client *Client
}

func (m *NamespaceMode) Init(client *Client) {
	m.client = client
}

func (m *NamespaceMode) Draw(v *View) error {
	for y, ns := range m.client.Namespaces.Items {
		x := 0
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault

		if y == v.ptr {
			bg = termbox.ColorGreen
		}

		for _, ch := range fmt.Sprintf("%02d", y) {
			termbox.SetCell(x, y, ch, fg, bg)
			x++
		}

		termbox.SetCell(x, y, ' ', fg, bg)
		x++

		for _, ch := range ns.ObjectMeta.Name {
			termbox.SetCell(x, y, ch, fg, bg)
			x++
		}

		for ; x < v.width; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, bg)
		}
	}

	return nil
}
