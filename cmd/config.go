package cmd

import (
	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
	"github.com/x-color/atchk/internal/atcoder"
)

func newConfigCmd() *cobra.Command {
	var langMode bool
	var conf *internal.Config
	at := atcoder.NewAtcoder()
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "edit & show config",
		Example: "",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			tmp, cache, err := internal.NewConfAndCache()
			if err != nil {
				return err
			}
			conf = tmp
			at.SetCache(cache)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if langMode {
				langs, err := at.GetLangList()
				if err != nil {
					return err
				}
				for _, lang := range langs {
					cmd.Printf("- %s\n", lang)
				}
				return nil
			} else {
				cmd.Println(conf.String())
				return nil
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&langMode, "lang-list", false, "list usable language for answer code")

	return cmd
}
