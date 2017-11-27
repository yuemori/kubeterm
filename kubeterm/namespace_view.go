package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"log"
)

type NamespaceView struct {
	watcher watch.Interface
	items   []v1.Namespace
}

func NewNamespaceView() *NamespaceView {
	return &NamespaceView{
		items: []v1.Namespace{},
	}
}

func (v *NamespaceView) Open(a *App, gv *gocui.View) {
	gv.SelBgColor = gocui.ColorRed
	gv.SelFgColor = gocui.ColorGreen
	v.ptrInit(gv)
	a.SetViewKeybinding(v, 'q', ModNone, v.quit)
	a.SetViewKeybinding(v, 'j', ModNone, v.ptrDown)
	a.SetViewKeybinding(v, 'k', ModNone, v.ptrUp)

	watcher := a.client.WatchNamespace(func(nss *v1.NamespaceList) {
		a.Update(func() error { return v.update(gv, nss) })
	})

	v.watcher = watcher
}

func (v *NamespaceView) update(gv *gocui.View, nss *v1.NamespaceList) error {
	gv.Clear()
	v.printLine(gv, "Name", "Status", "CreationTimestamp")

	for _, ns := range nss.Items {
		v.printLine(gv, ns.ObjectMeta.Name, ns.Status.Phase, ns.ObjectMeta.CreationTimestamp.Time)
	}

	v.items = nss.Items

	return nil
}

func (v *NamespaceView) printLine(gv *gocui.View, a ...interface{}) {
	fmt.Fprintln(gv, fmt.Sprintf("%-60s\t%-10s\t%-20s", a...))
}

func (v *NamespaceView) Close() {
	v.watcher.Stop()
}

func (v *NamespaceView) Name() string {
	return "namespace"
}

func (v *NamespaceView) DisplayName() string {
	return "Namespaces"
}

func (v *NamespaceView) quit(a *App, gv *gocui.View) error {
	a.Update(func() error {
		gv.Highlight = false
		return nil
	})
	a.ReturnToMenu()

	return nil
}

func (v *NamespaceView) OnFocus(a *App, gv *gocui.View) {
	v.ptrInit(gv)

	a.Update(func() error {
		gv.Highlight = true
		return v.update(gv, &v1.NamespaceList{
			Items: v.items,
		})
	})
}

func (v *NamespaceView) ptrInit(gv *gocui.View) {
	if err := gv.SetCursor(0, 1); err != nil {
		log.Panicln(err)
	}
}

func (v *NamespaceView) ptrDown(a *App, gv *gocui.View) error {
	x, y := gv.Cursor()
	next := y + 1

	if next > len(v.items) {
		next = y
	}

	if err := gv.SetCursor(x, next); err != nil {
		return err
	}

	return nil
}

func (v *NamespaceView) ptrUp(a *App, gv *gocui.View) error {
	x, y := gv.Cursor()
	next := y - 1

	if next < 1 {
		next = 1
	}

	if err := gv.SetCursor(x, next); err != nil {
		return err
	}

	return nil
}
