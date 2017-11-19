package kubeterm

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type PodMode struct {
	client    *Client
	namespace string
}

func NewPodMode(client *Client, namespace string) *PodMode {
	client.Clear()

	return &PodMode{
		client:    client,
		namespace: namespace,
	}
}

func (m *PodMode) Draw(ptr int, width int) error {
	for y, ns := range m.client.Pods(m.namespace).Items {
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

func (m *PodMode) Next(ptr int) Mode {
	return NewPodMode(m.client, m.namespace)
}

func (m *PodMode) Prev() Mode {
	return NewNamespaceMode(m.client)
}
