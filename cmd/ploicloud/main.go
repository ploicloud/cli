package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ploicloud/cli/internal/auth"
	"github.com/ploicloud/cli/internal/client"
	"github.com/ploicloud/cli/internal/commands"
	"github.com/ploicloud/cli/internal/commands/custom"
	"github.com/ploicloud/cli/internal/config"
	"github.com/ploicloud/cli/internal/output"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRoot().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRoot() *cobra.Command {
	invoked := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
	if invoked != "ploicloud" && invoked != "pcctl" {
		invoked = "ploicloud"
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load config: %v\n", err)
		cfg = &config.Config{APIURL: config.DefaultAPIURL, Output: config.DefaultOutput}
	}

	c := client.New(cfg).WithRefresher(auth.DefaultRefresher{})

	root := &cobra.Command{
		Use:           invoked,
		Short:         "Command line interface for Ploi Cloud",
		Long:          "ploicloud (alias: pcctl) wraps the Ploi Cloud REST API. Run `" + invoked + " login` to get started.",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	root.PersistentFlags().StringVar(&cfg.APIURL, "api-url", cfg.APIURL, "Base URL of the Ploi Cloud API")
	root.PersistentFlags().BoolVar(&output.JSON, "json", false, "Output raw JSON instead of formatted text")

	root.AddCommand(custom.NewLoginCommand())
	root.AddCommand(custom.NewLogoutCommand())
	root.AddCommand(custom.NewWhoamiCommand(c))
	root.AddCommand(custom.NewVersionCommand())

	commands.Register(root, c)
	return root
}
