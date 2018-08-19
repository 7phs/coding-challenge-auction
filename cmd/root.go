package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	ApplicationName = "auction-bit-tracker"
	Version         = "0.1"
)

var (
	GitHash   string // should be uninitialized
	BuildTime string // should be uninitialized
)

var RootCmd = &cobra.Command{
	Use:   ApplicationName,
	Short: "Auction Bit Tracker server",
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ApplicationInfo())
	},
}

func ApplicationInfo() string {
	return "Auction Bid Tracker " + Version + " [" + GitHash + "] " + BuildTime
}
