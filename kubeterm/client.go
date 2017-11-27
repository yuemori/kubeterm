package kubeterm

import (
	"github.com/thoas/go-funk"
	"github.com/yuemori/kubeterm/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"log"
)

type Client struct {
	client v1core.CoreV1Interface

	NamespaceList *v1.NamespaceList
}

func NewClient(config *Config) *Client {
	clientConfig := kubernetes.NewClientConfig(config.KubeConfig, config.ContextName)
	clientset, err := kubernetes.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}

	c := &Client{
		client:        clientset.CoreV1(),
		NamespaceList: &v1.NamespaceList{},
	}

	return c
}

func (c *Client) WatchNamespace(handler func(nss *v1.NamespaceList)) watch.Interface {
	watcher, err := c.client.Namespaces().Watch(v1.ListOptions{Watch: true})

	if err != nil {
		log.Panicln(err)
	}

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				if e.Object == nil {
					continue
				}

				ns := e.Object.(*v1.Namespace)
				switch e.Type {
				case watch.Added:
					c.addNamespace(ns)
				case watch.Modified:
					c.updateNamespace(ns)
				case watch.Error:
					c.updateNamespace(ns)
				case watch.Deleted:
					c.deleteNamespace(ns)
				}

				handler(c.NamespaceList)
			}
		}
	}()

	return watcher
}

func (c *Client) addNamespace(ns *v1.Namespace) {
	c.NamespaceList.Items = append(c.NamespaceList.Items, *ns)
}

func (c *Client) updateNamespace(ns *v1.Namespace) {
	c.deleteNamespace(ns)
	c.addNamespace(ns)
}

func (c *Client) deleteNamespace(namespace *v1.Namespace) {
	f := func(ns v1.Namespace) bool { return namespace.ObjectMeta.UID != ns.ObjectMeta.UID }
	c.NamespaceList.Items = funk.Filter(c.NamespaceList.Items, f).([]v1.Namespace)
}
