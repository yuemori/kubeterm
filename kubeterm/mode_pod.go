package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
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

func (m *PodMode) Draw(v *gocui.View) {
	// v.SetKeybinding("namespace", 'j')
	fmt.Fprintln(v, fmt.Sprintf("%-60s\t%-10s\t%-20s", "Name", "Status", "CreationTimestamp"))

	pods := m.client.Pods(m.namespace).Items

	for _, p := range pods {
		fmt.Fprintln(v, fmt.Sprintf("%-60s\t%-10s\t%-20s", p.ObjectMeta.Name, p.Status.Phase, p.ObjectMeta.CreationTimestamp.Time))
	}
}
