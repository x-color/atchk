package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
	"github.com/x-color/atchk/internal/atcoder"
)

func newLoginCmd() *cobra.Command {
	at := atcoder.NewAtcoder()
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "login atcoder",
		Example: "",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return at.LoadConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if at.IsLoggedIn() {
				return nil
			}

			var username string
			cmd.Print("Enter your user name: ")
			fmt.Scan(&username)
			cmd.Print("Enter your password : ")
			password, err := internal.ReadPwd()
			cmd.Println("")
			if err != nil {
				return err
			}

			if err := at.Login(username, password); err != nil {
				return err
			}

			return at.SaveConfig()
		},
	}

	return cmd
}
