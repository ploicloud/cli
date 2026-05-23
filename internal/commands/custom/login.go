package custom

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ploicloud/cli/internal/auth"
	"github.com/spf13/cobra"
)

func NewLoginCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Ploi Cloud via your browser",
		Long:  "Opens a browser window to authorize the CLI, then stores the resulting OAuth token in your OS keychain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if auth.ClientID == "" {
				return fmt.Errorf("CLI OAuth client_id is not configured; rebuild with `make build CLIENT_ID=<uuid>`")
			}

			pkce, err := auth.NewPKCE()
			if err != nil {
				return err
			}

			lb, err := auth.StartLoopback(pkce.State)
			if err != nil {
				return err
			}
			defer lb.Close()

			authURL, err := auth.BuildAuthorizeURL(auth.AuthorizeURL, auth.ClientID, lb.RedirectURI(), auth.Scope, pkce)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Opening browser to authorize.\nIf nothing happens, visit:\n  %s\n\nWaiting for the callback on %s ...\n", authURL, lb.RedirectURI())
			if berr := auth.OpenBrowser(authURL); berr != nil {
				fmt.Fprintf(os.Stderr, "(could not open browser automatically: %v)\n", berr)
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
			defer cancel()

			res, err := lb.Wait(ctx)
			if err != nil {
				return err
			}

			tok, err := auth.ExchangeCode(ctx, res.Code, pkce.Verifier, lb.RedirectURI())
			if err != nil {
				return err
			}
			if err := auth.Save(tok); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "Logged in.")
			return nil
		},
	}
}
