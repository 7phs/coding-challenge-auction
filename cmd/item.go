package cmd

import (
	"fmt"
	"os"

	"github.com/7phs/coding-challenge-auction/api"
	"github.com/7phs/coding-challenge-auction/models"
	"github.com/spf13/cobra"
)

var (
	addItemTitle *string

	updateItemId    *int
	updateItemTitle *string

	getItemId *int
	topItemId *int
)

func init() {
	addItemTitle = ItemAddCmd.Flags().String("title", "", "an item's title")

	updateItemId = ItemUpdateCmd.Flags().Int("id", 0, "an item's id")
	updateItemTitle = ItemUpdateCmd.Flags().String("title", "", "an item's title")

	getItemId = ItemGetCmd.Flags().Int("id", 0, "an item's id")
	topItemId = ItemTopCmd.Flags().Int("id", 0, "an item's id")
}

var ItemCmd = &cobra.Command{
	Use:   "item",
	Short: "Commands processing an item",
}

var ItemAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an item",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := api.NewItem().Add(*addItemTitle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			fmt.Printf("Success: #%d\n", key)
		}
	},
}

var ItemUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an item",
	Run: func(cmd *cobra.Command, args []string) {
		err := api.NewItem().Update(models.ItemKey(*updateItemId), *updateItemTitle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			fmt.Printf("Success\n")
		}
	},
}

var ItemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an item and all bids",
	Run: func(cmd *cobra.Command, args []string) {
		item, err := api.NewItem().Get(models.ItemKey(*getItemId))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		fmt.Printf("Id:    %d\n", item.Data.Id)
		fmt.Printf("Title: %s\n", item.Data.Title)
		if len(item.Data.Bids) == 0 {
			fmt.Println("Bids: no one")
		} else {
			fmt.Println("Bids:")
			for _, bid := range item.Data.Bids {
				fmt.Printf("    #%d '%s' %s %s\n", bid.UserId, bid.UserName, bid.Bid, bid.Updated)
			}
		}
	},
}

var ItemTopCmd = &cobra.Command{
	Use:   "top",
	Short: "Get the top item's bid",
	Run: func(cmd *cobra.Command, args []string) {
		top, err := api.NewItem().Top(models.ItemKey(*topItemId))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		fmt.Printf("Top bid: user #%d '%s' %s %s\n", top.Data.UserId, top.Data.UserName, top.Data.Bid, top.Data.Updated)
	},
}
