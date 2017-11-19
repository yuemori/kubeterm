package kubeterm

import (
	"context"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
)

func stdErr(err error) {
	fmt.Fprintf(os.Stderr, "kome: %v\n", err)
}

func Run(ctx context.Context, config *Config) error {
	// init termbox
	if err := termbox.Init(); err != nil {
		stdErr(err)
		return err
	}
	defer termbox.Close()

	client := NewClient(config)
	view := NewView(client)
	view.Loop()

	return nil
}
