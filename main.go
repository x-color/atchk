package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  %s CONTEST COMMAND\n", os.Args[0])
		fmt.Println("e.g.:")
		fmt.Printf("  %s abc001a go run main.go\n", os.Args[0])
		fmt.Printf("  %s abc001a python main.py\n", os.Args[0])
		os.Exit(1)
	}
	execChecker(os.Args[1], os.Args[2:])
}

func execChecker(name string, command []string) {
	fmt.Printf("Check %s\n\n", strings.ToUpper(name))
	contest, err := getContest(strings.ToLower(name))
	if err != nil {
		fmt.Println(err)
		return
	}

	result := checkCode(command, contest)

	fmt.Println()
	if result {
		fmt.Println("Your code is OK!!!")
	} else {
		fmt.Println("Your code is NG...")
	}
}
