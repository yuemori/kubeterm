package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/thoas/go-funk"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"log"
)

type PodView struct {
	watcher   watch.Interface
	client    *Client
	namespace string
	items     []v1.Pod
}

func NewPodView(ns string, client *Client) *PodView {
	return &PodView{
		namespace: ns,
		client:    client,
		items:     []v1.Pod{},
	}
}

func (v *PodView) Open(a *App, gv *gocui.View) {
	gv.SelBgColor = gocui.ColorRed
	gv.SelFgColor = gocui.ColorGreen
	a.SetViewKeybinding(v, 'q', ModNone, v.quit)
	a.SetViewKeybinding(v, 'j', ModNone, v.ptrDown)
	a.SetViewKeybinding(v, 'k', ModNone, v.ptrUp)

	v.init(a.CurrentNamespace(), a, gv)
}

func (v *PodView) quit(a *App, gv *gocui.View) error {
	a.Update(func() error {
		gv.Highlight = false
		return nil
	})
	a.ReturnToMenu()

	return nil
}

func (v *PodView) ptrInit(gv *gocui.View) {
	if err := gv.SetCursor(0, 1); err != nil {
		log.Panicln(err)
	}
}

func (v *PodView) update(gv *gocui.View, pods *v1.PodList) error {
	gv.Clear()
	v.printLine(gv, "Name", "Ready", "Status", "CreationTimestamp")

	for _, p := range pods.Items {
		max := len(p.Status.ContainerStatuses)
		now := len(funk.Filter(p.Status.ContainerStatuses, func(s v1.ContainerStatus) bool { return s.Ready }).([]v1.ContainerStatus))

		ready := fmt.Sprintf("%d/%d", now, max)

		v.printLine(gv, p.ObjectMeta.Name, ready, p.Status.Phase, p.ObjectMeta.CreationTimestamp.Time)
	}

	v.items = pods.Items

	return nil
}

func (v *PodView) printLine(gv *gocui.View, a ...interface{}) {
	fmt.Fprintln(gv, fmt.Sprintf("%-60s\t%-10s\t%-10s\t%-20s", a...))
}

func (v *PodView) Close() {
	v.watcher.Stop()
}

func (v *PodView) init(ns string, a *App, gv *gocui.View) {
	v.namespace = ns

	if v.watcher != nil {
		v.watcher.Stop()
	}

	v.items = []v1.Pod{}

	v.update(gv, &v1.PodList{Items: v.items})

	watcher := v.client.WatchPod(ns, func(pods *v1.PodList) {
		a.Update(func() error {
			return v.update(gv, pods)
		})
	})

	v.watcher = watcher
}

func (v *PodView) OnEnter(a *App, gv *gocui.View) {
	gv.Highlight = true
	v.OnFocus(a, gv)
}

func (v *PodView) OnFocus(a *App, gv *gocui.View) {
	v.ptrInit(gv)

	if ns := a.CurrentNamespace(); ns != v.namespace {
		v.init(ns, a, gv)
		return
	}

	a.Update(func() error {
		return v.update(gv, &v1.PodList{
			Items: v.items,
		})
	})
}

func (v *PodView) Name() string {
	return "pod"
}

func (v *PodView) DisplayName() string {
	return "Pods"
}

func (v *PodView) ptrDown(a *App, gv *gocui.View) error {
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

func (v *PodView) ptrUp(a *App, gv *gocui.View) error {
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
