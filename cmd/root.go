package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
)

var config internal.Config

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "atchk",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newLogoutCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newEnterCmd())

	return cmd
}

func Execute() {
	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}
