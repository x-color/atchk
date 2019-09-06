package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
	"github.com/x-color/atchk/internal/atcoder"
)

func newSubmitCmd() *cobra.Command {
	var conf *internal.Config
	at := atcoder.NewAtcoder()
	cmd := &cobra.Command{
		Use:     "submit",
		Short:   "submit answer code",
		Example: "",
		Args:    cobra.ExactArgs(2),
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
			contest := args[0]
			if len(contest) != 7 || (contest[:3] != "abc" && contest[:3] != "agc") {
				return fmt.Errorf("Invalid contest name")
			}
			if _, err := strconv.Atoi(contest[3:6]); err != nil {
				return fmt.Errorf("Invalid contest name")
			}
			return at.Submit(args[0][:6], args[0][6:7], args[1], conf.LangID)
		},
	}

	return cmd
}
