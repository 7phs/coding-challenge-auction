package main

import "github.com/7phs/coding-challenge-auction/cmd"

func main() {
	cmd.RootCmd.AddCommand(
		cmd.VersionCmd,
		cmd.RunCmd)
}
