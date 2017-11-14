package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yuemori/kubeterm/kubeterm"
)

const version = "0.0.1"

type Options struct {
	kubeConfig string
	version    bool
	context    string
}

var opts = &Options{}

func Run() {
	cmd := &cobra.Command{}
	cmd.Use = "kubeterm"
	cmd.Short = "Kubernetes interactive terminal."

	cmd.Flags().StringVar(&opts.context, "context", opts.context, "Kubernetes context to use. Default to current context configured in kubeconfig.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if opts.version {
			fmt.Printf("kubeterm version %s\n", version)
			return nil
		}

		narg := len(args)
		if narg > 1 {
			return cmd.Help()
		}
		config, err := parseConfig(args)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err = kubeterm.Run(ctx, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(args []string) (*kubeterm.Config, error) {
	kubeConfig, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	return &kubeterm.Config{
		KubeConfig:  kubeConfig,
		ContextName: opts.context,
	}, nil
}

func getKubeConfig() (string, error) {
	var kubeconfig string

	if kubeconfig = opts.kubeConfig; kubeconfig != "" {
		return kubeconfig, nil
	}

	if kubeconfig = os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig, nil
	}

	// kubernetes requires an absolute path
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home directory")
	}

	kubeconfig = filepath.Join(home, ".kube/config")

	return kubeconfig, nil
}
