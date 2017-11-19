package kubeterm

import (
	"github.com/yuemori/kubeterm/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

type Client struct {
	client v1core.CoreV1Interface
	v1core.CoreV1Interface
	Namespaces *v1.NamespaceList
}

func NewClient(config *Config) *Client {
	clientConfig := kubernetes.NewClientConfig(config.KubeConfig, config.ContextName)
	clientset, err := kubernetes.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}

	client := clientset.CoreV1()

	nss, _ := client.Namespaces().List(v1.ListOptions{})

	c := &Client{
		client:     clientset.CoreV1(),
		Namespaces: nss,
	}

	c.watchNamespace()

	return c
}

func (c *Client) watchNamespace() error {
	watcher, err := c.client.Namespaces().Watch(v1.ListOptions{Watch: true})

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				switch e.Type {
				case watch.Added:
					nss, _ := c.client.Namespaces().List(v1.ListOptions{})
					c.Namespaces = nss
				case watch.Deleted:
					nss, _ := c.client.Namespaces().List(v1.ListOptions{})
					c.Namespaces = nss
				}
			}
		}
	}()

	return nil
}
