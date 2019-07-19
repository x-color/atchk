package main

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func checkCode(command []string, contest *Contest) bool {
	r := true
	for _, samples := range contest.Samples {
		res, err := execSample(command, samples)
		if err != nil {
			fmt.Println(colorRed("\nError Occured!!:"))
			fmt.Println(err)
			return false
		}
		r = r && res
	}
	return r
}

func getCmd(command []string, input string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	if len(command) == 1 {
		cmd = exec.Command(command[0])
	} else {
		cmd = exec.Command(command[0], command[1:]...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()
	_, err = io.WriteString(stdin, input)

	return cmd, err
}

func execSample(command []string, sample *Sample) (bool, error) {
	cmd, err := getCmd(command, sample.Input)
	if err != nil {
		return false, err
	}

	fmt.Printf("- Execute %s ... ", sample.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

	return checkResult(string(output), sample.Output), nil
}

func checkResult(act, exp string) bool {
	if act != exp {
		fmt.Println(colorRed("FAILURE"))
		fmt.Printf("    expected:%s", exp)
		fmt.Printf("    actual  :%s", act)
		return false
	}
	fmt.Println(colorGreen("SUCCESS"))
	return true
}

func colorRed(txt string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", txt)
}

func colorGreen(txt string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", txt)
}
