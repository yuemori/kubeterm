package kubeterm

import (
	"context"
	"fmt"
	"os"
)

func stdErr(err error) {
	fmt.Fprintf(os.Stderr, "kome: %v\n", err)
}

func Run(ctx context.Context, config *Config) error {
	client := NewClient(config)
	app := NewApp(client)
	app.MainLoop()

	return nil
}
