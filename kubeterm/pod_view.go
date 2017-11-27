package kubeterm

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/thoas/go-funk"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

type PodView struct {
	watcher   watch.Interface
	namespace string
}

func NewPodView(ns string) *PodView {
	return &PodView{
		namespace: ns,
	}
}

func (v *PodView) Open(a *App, gv *gocui.View) {
	gv.SelBgColor = gocui.ColorRed
	gv.SelFgColor = gocui.ColorGreen

	watcher := a.client.WatchPod(v.namespace, func(pods *v1.PodList) {
		a.g.Update(func(g *gocui.Gui) error {
			gv.Clear()
			v.printLine(gv, "Name", "Ready", "Status", "CreationTimestamp")

			for _, p := range pods.Items {
				max := len(p.Status.ContainerStatuses)
				now := len(funk.Filter(p.Status.ContainerStatuses, func(s v1.ContainerStatus) bool { return s.Ready }).([]v1.ContainerStatus))

				ready := fmt.Sprintf("%d/%d", now, max)

				v.printLine(gv, p.ObjectMeta.Name, ready, p.Status.Phase, p.ObjectMeta.CreationTimestamp.Time)
			}

			return nil
		})
	})

	v.watcher = watcher
}

func (v *PodView) printLine(gv *gocui.View, a ...interface{}) {
	fmt.Fprintln(gv, fmt.Sprintf("%-60s\t%-10s\t%-10s\t%-20s", a...))
}

func (v *PodView) Close() {
	v.watcher.Stop()
}

func (v *PodView) OnEnter() {
}

func (v *PodView) Name() string {
	return "pod"
}

func (v *PodView) DisplayName() string {
	return "Pods"
}
