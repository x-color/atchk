package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func ReadPwd() (string, error) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	defer signal.Stop(signalChan)

	originalTerminalState, err := terminal.GetState(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	go func() {
		<-signalChan
		terminal.Restore(int(syscall.Stdin), originalTerminalState)
		fmt.Println()
		os.Exit(1)
	}()

	password, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	return string(password), nil
}
