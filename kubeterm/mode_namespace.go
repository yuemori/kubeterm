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

func (m *NamespaceMode) Draw() error {
	bg := termbox.ColorDefault

	for y, ns := range m.client.Namespaces.Items {
		x := 0
		for _, ch := range fmt.Sprintf("%02d", y) {
			termbox.SetCell(x, y, ch, termbox.ColorBlue, bg)
			x++
		}

		termbox.SetCell(x, y, ' ', termbox.ColorBlue, bg)
		x++

		for _, ch := range ns.ObjectMeta.Name {
			termbox.SetCell(x, y, ch, termbox.ColorBlue, bg)
			x++
		}
	}

	return nil
}
