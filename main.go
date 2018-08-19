package main

import (
	"fmt"
	"os"

	"github.com/7phs/coding-challenge-auction/cmd"
)

func main() {
	cmd.RootCmd.AddCommand(
		cmd.VersionCmd,
		cmd.RunCmd,
		cmd.BidCmd,
		cmd.ItemCmd,
		cmd.UserCmd)

	cmd.BidCmd.AddCommand(
		cmd.BidPushCmd)

	cmd.ItemCmd.AddCommand(
		cmd.ItemAddCmd,
		cmd.ItemUpdateCmd,
		cmd.ItemGetCmd,
		cmd.ItemTopCmd)

	cmd.UserCmd.AddCommand(
		cmd.UserAddCmd,
		cmd.UserUpdateCmd,
		cmd.UserGetCmd)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
