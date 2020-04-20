package main

import (
	"os"

	"github.com/falmotlag/aws-bulk-cli/cli"
	"github.com/fatih/color"
)

func main() {
	err := cli.Run()
	checkForErrorsAndExit(err)
}

func checkForErrorsAndExit(err error) {
	if err == nil {
		os.Exit(0)
	} else {
		color.New(color.FgRed).Println(err.Error())
		os.Exit(1)
	}
}
