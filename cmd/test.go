package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/x-color/atchk/internal"
	"github.com/x-color/atchk/internal/atcoder"
)

func newTestCmd() *cobra.Command {
	var conf *internal.Config
	at := atcoder.NewAtcoder()
	cmd := &cobra.Command{
		Use:     "test",
		Short:   "test the answer code",
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
			n, err := at.Samples(args[0][:6], args[0][6:7])
			if err != nil {
				return err
			}
			cmdStr := strings.TrimSpace(fmt.Sprintf("%s %s", conf.Command, args[1]))
			for i := 0; i < n; i++ {
				msg, err := at.Test(i, strings.Split(cmdStr, " "))
				if err != nil {
					return err
				}
				cmd.Println(msg)
			}
			return nil
		},
	}

	return cmd
}
