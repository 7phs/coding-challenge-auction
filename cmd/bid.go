package cmd

import (
	"fmt"
	"os"

	"github.com/7phs/coding-challenge-auction/api"
	"github.com/7phs/coding-challenge-auction/models"
	"github.com/spf13/cobra"
)

var (
	bidItemId *int
	bidUserId *int
	bid       *float64
)

func init() {
	bidItemId = BidPushCmd.Flags().Int("item_id", 0, "an item id")
	bidUserId = BidPushCmd.Flags().Int("user_id", 0, "a user id")
	bid = BidPushCmd.Flags().Float64("bid", 0, "a bid")
}

var BidCmd = &cobra.Command{
	Use:   "bid",
	Short: "Commands processing a bid",
}

var BidPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a user's bid for an item",
	Run: func(cmd *cobra.Command, args []string) {
		err := api.NewBid(models.ItemKey(*bidItemId), models.UserKey(*bidUserId)).Push(*bid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			fmt.Println("Success")
		}
	},
}
