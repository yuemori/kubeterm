package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type NamespaceMode struct {
	client *Client
}

func NewNamespaceMode(client *Client) *NamespaceMode {
	client.Clear()

	return &NamespaceMode{
		client: client,
	}
}

func (m *NamespaceMode) Draw(v *gocui.View) {
	fmt.Fprintln(v, fmt.Sprintf("%-20s\t%-10s\t%-20s", "Name", "Status", "CreationTimestamp"))

	nss := m.client.Namespaces().Items

	for _, ns := range nss {
		fmt.Fprintln(v, fmt.Sprintf("%-20s\t%-10s\t%-20s", ns.ObjectMeta.Name, ns.Status.Phase, ns.ObjectMeta.CreationTimestamp.Time))
	}
}
