package contest

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/logrusorgru/aurora"
)

type Samples []*Sample

type Sample struct {
	ID     int    `mapstructure:"id"`
	Input  string `mapstructure:"input"`
	Output string `mapstructure:"output"`
}

func (sample *Sample) Test(cmdList []string) (string, error) {
	var cmd *exec.Cmd
	if len(cmdList) == 1 {
		cmd = exec.Command(cmdList[0])
	} else {
		cmd = exec.Command(cmdList[0], cmdList[1:]...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(stdin, sample.Input)
	stdin.Close()

	output, err := sample.exec(cmd)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("- Execute Sample %d --> %s",
		sample.ID, sample.check(output, sample.Output)), nil
}

func (sample *Sample) exec(cmd *exec.Cmd) (string, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) == 0 {
			return "", err
		}
		return "", fmt.Errorf("%s", output)
	}
	return string(output), nil
}

func (sample *Sample) check(act, exp string) string {
	if act != exp {
		return fmt.Sprintf("%s\n[*] expected:\n%s\n[*] actual:\n%s",
			aurora.Bold(aurora.Red("FAILURE")), exp, act)
	}
	return fmt.Sprint(aurora.Bold(aurora.Green("SUCCESS")))
}
