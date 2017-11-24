package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type NamespaceView struct {
	done chan struct{}
}

func NewNamespaceView() *NamespaceView {
	return &NamespaceView{
		done: make(chan struct{}),
	}
}

func (v *NamespaceView) Open(a *App, gv *gocui.View) {
	gv.SelBgColor = gocui.ColorRed
	gv.SelFgColor = gocui.ColorGreen

	a.client.WatchNamespace(func(nss *v1.NamespaceList) {
		v.draw(a, gv, nss)
	})
}

func (v *NamespaceView) draw(a *App, gv *gocui.View, nss *v1.NamespaceList) {
	gv.Clear()

	v.printLine(gv, "Name", "Status", "CreationTimestamp")

	for _, ns := range nss.Items {
		v.printLine(gv, ns.ObjectMeta.Name, ns.Status.Phase, ns.ObjectMeta.CreationTimestamp.Time)
	}
}

func (v *NamespaceView) printLine(gv *gocui.View, a ...interface{}) {
	fmt.Fprintln(gv, fmt.Sprintf("%-60s\t%-10s\t%-20s", a...))
}

func (v *NamespaceView) Close() {
	close(v.done)
}

func (v *NamespaceView) Name() string {
	return "detail"
}
