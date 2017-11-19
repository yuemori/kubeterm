package kubeterm

import (
	"github.com/yuemori/kubeterm/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client struct {
	client v1core.CoreV1Interface
	v1core.CoreV1Interface
}

func NewClient(config *Config) *Client {
	clientConfig := kubernetes.NewClientConfig(config.KubeConfig, config.ContextName)
	clientset, err := kubernetes.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}

	return &Client{
		client: clientset.CoreV1(),
	}
}

func (c *Client) Namespaces() v1.NamespaceInterface {
	return c.client.Namespaces()
}

func (c *Client) Pods(namespace string) v1.PodInterface {
	return c.client.Pods(namespace)
}
