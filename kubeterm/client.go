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
	namespaces *v1.NamespaceList
	pods       *v1.PodList
}

func NewClient(config *Config) *Client {
	clientConfig := kubernetes.NewClientConfig(config.KubeConfig, config.ContextName)
	clientset, err := kubernetes.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}

	c := &Client{
		client: clientset.CoreV1(),
	}

	return c
}

func (c *Client) Namespaces() *v1.NamespaceList {
	if c.namespaces == nil {
		nss, _ := c.client.Namespaces().List(v1.ListOptions{})
		c.namespaces = nss
		if err := c.watchNamespace(); err != nil {
			panic(err)
		}
	}

	return c.namespaces
}

func (c *Client) Pods(namespace string) *v1.PodList {
	if c.pods == nil {
		pods, _ := c.client.Pods(namespace).List(v1.ListOptions{})
		c.pods = pods
		if err := c.watchPod(namespace); err != nil {
			panic(err)
		}
	}

	return c.pods
}

func (c *Client) watchPod(namespace string) error {
	watcher, err := c.client.Pods(namespace).Watch(v1.ListOptions{Watch: true})

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case e := <-watcher.ResultChan():
				switch e.Type {
				case watch.Added:
					pods, _ := c.client.Pods(namespace).List(v1.ListOptions{})
					c.pods = pods
				case watch.Deleted:
					pods, _ := c.client.Pods(namespace).List(v1.ListOptions{})
					c.pods = pods
				}
			}
		}
	}()

	return nil
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
					c.namespaces = nss
				case watch.Deleted:
					nss, _ := c.client.Namespaces().List(v1.ListOptions{})
					c.namespaces = nss
				}
			}
		}
	}()

	return nil
}
