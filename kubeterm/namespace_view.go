package kubeterm

import (
	v1 "k8s.io/client-go/pkg/api/v1"
	"log"
)

type NamespaceView struct {
	*TableView
	client *Client
	items  []v1.Namespace
}

func NewNamespaceView(client *Client) *NamespaceView {
	table := NewTableView("%-60s\t%-10s\t%-20s")
	table.AddHeader("Name", "Status", "CreationTimestamp")

	nss, err := client.Interface.Namespaces().List(v1.ListOptions{})

	if err != nil {
		log.Panicln(err)
	}

	for _, ns := range nss.Items {
		table.AddRow(ns.ObjectMeta.Name, ns.Status.Phase, ns.ObjectMeta.CreationTimestamp.Time)
	}

	return &NamespaceView{
		TableView: table,
		client:    client,
		items:     nss.Items,
	}
}

func (v *NamespaceView) Init(view *View) {
	view.SetKeybinding('q', view.Quit)
	view.SetKeybinding('k', view.PointerUp)
	view.SetKeybinding('j', view.PointerDown)
	// view.SetKeybinding(KeyEnter, func() error {
	// 	return nil
	// })
}

func (v *NamespaceView) BeginPointerIndex() (x int) {
	return len(v.Headers)
}

func (v *NamespaceView) Position() (x0, y0, x1, y1 int) {
	w := GetWindow()
	return 20, 0, w.Width, w.Height
}

func (v *NamespaceView) Name() string {
	return "Namespaces"
}
