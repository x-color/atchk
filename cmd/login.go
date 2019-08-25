package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
	"github.com/x-color/atchk/internal/atcoder"
)

func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "login atcoder",
		Example: "",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.Read()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var at atcoder.Atcoder
			if at.IsLoggedIn(config.System.Cookies) {
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

			cookies, err := at.Login(username, password)
			if err != nil {
				return err
			}
			config.System.Cookies = cookies

			return config.Update()
		},
	}

	return cmd
}
