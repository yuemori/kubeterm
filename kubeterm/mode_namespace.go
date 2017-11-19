package kubeterm

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type NamespaceMode struct {
	client *Client
}

func NewNamespaceMode(client *Client) *NamespaceMode {
	return &NamespaceMode{
		client: client,
	}
}

func (m *NamespaceMode) Draw(ptr int, width int) error {
	for y, ns := range m.client.Namespaces().Items {
		x := 0
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault

		if y == ptr {
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

		for ; x < width; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, bg)
		}
	}

	return nil
}

func (m *NamespaceMode) Next(ptr int) Mode {
	return NewPodMode(m.client, m.current(ptr))
}

func (m *NamespaceMode) Prev() Mode {
	return nil
}

func (m *NamespaceMode) current(ptr int) (namespace string) {
	nss := m.client.Namespaces().Items
	return nss[ptr].Name
}
