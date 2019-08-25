package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func newEnterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "enter",
		Short:   "enter the contest",
		Example: "",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.Read()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			contest := strings.ToLower(args[0])
			if err := config.ValidContest(contest); err != nil {
				return err
			}
			config.System.Contest = contest
			return config.Update()
		},
	}

	return cmd
}
