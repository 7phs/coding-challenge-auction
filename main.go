package main

import (
	"fmt"
	"os"

	"github.com/7phs/coding-challenge-auction/cmd"
)

func main() {
	cmd.RootCmd.AddCommand(
		cmd.VersionCmd,
		cmd.RunCmd)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
