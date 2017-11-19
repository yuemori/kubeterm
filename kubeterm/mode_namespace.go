package kubeterm

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"k8s.io/client-go/pkg/api/v1"
)

type NamespaceMode struct {
	client *Client
}

func (m *NamespaceMode) Init(client *Client) {
	m.client = client
}

func (m *NamespaceMode) Draw() error {
	nss, err := m.client.Namespaces().List(v1.ListOptions{})
	if err != nil {
		return err
	}

	bg := termbox.ColorDefault

	for y, ns := range nss.Items {
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
