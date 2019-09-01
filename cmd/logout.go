package cmd

import (
	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
	"github.com/x-color/atchk/internal/atcoder"
)

func newLogoutCmd() *cobra.Command {
	at := atcoder.NewAtcoder()
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "logout atcoder",
		Example: "",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, cache, err := internal.NewConfAndCache()
			if err != nil {
				return err
			}
			at.SetCache(cache)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			at.Logout()
			return at.SaveCache()
		},
	}

	return cmd
}
