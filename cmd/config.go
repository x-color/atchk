package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "edit & show config",
		Example: "",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 && len(args) != 2 {
				return fmt.Errorf("accepts 0 or 2 args, received %d", len(args))
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.Read()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				s, err := config.Format()
				if err != nil {
					return err
				}
				cmd.Println(s)
				return nil
			}
			if err := config.Set(args[0], args[1]); err != nil {
				return err
			}
			return config.Update()
		},
	}

	return cmd
}
