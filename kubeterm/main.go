package kubeterm

import (
	"context"
	"fmt"
)

func Run(ctx context.Context, config *Config) error {
	fmt.Printf("hello, kubeterm\n")
	return nil
}
