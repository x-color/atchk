package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal/atcoder"
)

func newEnterCmd() *cobra.Command {
	at := atcoder.NewAtcoder()
	cmd := &cobra.Command{
		Use:     "enter",
		Short:   "enter the contest",
		Example: "",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return at.LoadConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := at.SetConfig("system.contest", strings.ToLower(args[0])); err != nil {
				return err
			}
			return at.SaveConfig()
		},
	}

	return cmd
}
