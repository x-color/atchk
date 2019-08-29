package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal/atcoder"
)

func newConfigCmd() *cobra.Command {
	var langMode bool
	at := atcoder.NewAtcoder()
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
			return at.LoadConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				if langMode {
					langs, err := at.GetLangList()
					if err != nil {
						return err
					}
					for i, lang := range langs {
						cmd.Printf("%2d %s\n", i+1, lang)
					}

					var i int
					cmd.Print("\nPlease select using language (enter the number):")
					fmt.Scanf("%d", &i)
					if i <= 0 || i > len(langs)  {
						return fmt.Errorf("Please select number 1 to %d", len(langs))
					}
					if err := at.SetLang(langs[i-1]); err != nil {
						return err
					}
				} else {
					cmd.Println(at.String())
					return nil
				}
			} else {
				if err := at.SetConfig(args[0], args[1]); err != nil {
					return err
				}
			}
			return at.SaveConfig()
		},
	}

	cmd.Flags().BoolVar(&langMode, "set-lang", false, "set using language for answer code")

	return cmd
}
