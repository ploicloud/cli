package custom

import (
	"context"

	"github.com/ploicloud/cli/internal/client"
	"github.com/ploicloud/cli/internal/output"
	"github.com/spf13/cobra"
)

func NewWhoamiCommand(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show the currently authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := c.Do(context.Background(), client.Request{
				Method:   "GET",
				PathTmpl: "/user",
			})
			if err != nil {
				return err
			}
			return output.Print(resp)
		},
	}
}
