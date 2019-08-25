package cmd

import (
	"github.com/spf13/cobra"
)

func newLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "logout atcoder",
		Example: "",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.Read()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config.System.Cookies = nil
			return config.Update()
		},
	}

	return cmd
}
