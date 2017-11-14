package kubeterm

import (
	"context"
	"fmt"

	"github.com/yuemori/kubeterm/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

func Run(ctx context.Context, config *Config) error {
	clientConfig := kubernetes.NewClientConfig(config.KubeConfig, config.ContextName)
	clientset, err := kubernetes.NewClientSet(clientConfig)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods("").List(v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, p := range pods.Items {
		fmt.Printf("Name: %s\n", p.ObjectMeta.Name)
	}

	return nil
}
