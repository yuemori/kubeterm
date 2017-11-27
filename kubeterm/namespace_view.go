package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

type NamespaceView struct {
	done    chan struct{}
	watcher watch.Interface
}

func NewNamespaceView() *NamespaceView {
	return &NamespaceView{
		done: make(chan struct{}),
	}
}

func (v *NamespaceView) Open(a *App, gv *gocui.View) {
	gv.SelBgColor = gocui.ColorRed
	gv.SelFgColor = gocui.ColorGreen

	watcher := a.client.WatchNamespace(func(nss *v1.NamespaceList) {
		a.g.Update(func(g *gocui.Gui) error {
			gv.Clear()
			v.printLine(gv, "Name", "Status", "CreationTimestamp")

			for _, ns := range nss.Items {
				v.printLine(gv, ns.ObjectMeta.Name, ns.Status.Phase, ns.ObjectMeta.CreationTimestamp.Time)
			}

			return nil
		})
	})

	v.watcher = watcher
}

func (v *NamespaceView) printLine(gv *gocui.View, a ...interface{}) {
	fmt.Fprintln(gv, fmt.Sprintf("%-60s\t%-10s\t%-20s", a...))
}

func (v *NamespaceView) Close() {
	v.watcher.Stop()
	close(v.done)
}

func (v *NamespaceView) Name() string {
	return "namespace"
}
