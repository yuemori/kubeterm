package kubeterm

import (
	"github.com/yuemori/kubeterm/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
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

	return c
}
