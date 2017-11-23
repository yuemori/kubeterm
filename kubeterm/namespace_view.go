package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
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

	v.draw(a, gv)
}

func (v *NamespaceView) draw(a *App, gv *gocui.View) {
	v.printLine(gv, "Name", "Status", "CreationTimestamp")

	nss := a.client.Namespaces()

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
