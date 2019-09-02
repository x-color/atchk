package cmd

import (
	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "initialize config of atchk",
		Example: "",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.CreateNewConfFile()
		},
	}

	return cmd
}
