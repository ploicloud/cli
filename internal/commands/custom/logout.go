package custom

import (
	"fmt"

	"github.com/ploicloud/cli/internal/auth"
	"github.com/spf13/cobra"
)

func NewLogoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.Clear(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStderr(), "Logged out.")
			return nil
		},
	}
}
