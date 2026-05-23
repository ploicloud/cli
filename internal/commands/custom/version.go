package custom

import (
	"fmt"

	"github.com/ploicloud/cli/internal/client"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "ploicloud %s (commit %s)\n", client.Version, client.Commit)
		},
	}
}
