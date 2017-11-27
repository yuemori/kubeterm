package kubeterm

import (
	"fmt"
	"github.com/thoas/go-funk"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"log"
)

type PodView struct {
	*TableView
	watcher   watch.Interface
	client    *Client
	namespace string
	items     []v1.Pod
}

func NewPodView(ns string, client *Client) *PodView {
	table := NewTableView("%-60s\t%-10s\t%-10s\t%-20s")
	table.AddHeader("Name", "Ready", "Status", "CreationTimestamp")
	pods, err := client.Interface.Pods(ns).List(v1.ListOptions{})

	if err != nil {
		log.Panicln(err)
	}

	for _, pod := range pods.Items {
		max := len(pod.Status.ContainerStatuses)
		now := len(funk.Filter(pod.Status.ContainerStatuses, func(s v1.ContainerStatus) bool { return s.Ready }).([]v1.ContainerStatus))
		ready := fmt.Sprintf("%d/%d", now, max)

		table.AddRow(pod.ObjectMeta.Name, ready, pod.Status.Phase, pod.ObjectMeta.CreationTimestamp.Time)
	}

	return &PodView{
		TableView: table,
		namespace: ns,
		client:    client,
		items:     pods.Items,
	}
}

func (v *PodView) Init(view *View) {
	view.SetKeybinding('q', view.Quit)
	view.SetKeybinding('k', view.PointerUp)
	view.SetKeybinding('j', view.PointerDown)

}

func (v *PodView) BeginPointerIndex() (x int) {
	return len(v.Headers)
}

func (v *PodView) Position() (x0, y0, x1, y1 int) {
	w := GetWindow()
	return 20, 0, w.Width, w.Height
}

func (v *PodView) Name() string {
	return "pod"
}

func (v *PodView) DisplayName() string {
	return "Pods"
}
