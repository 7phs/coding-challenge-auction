package cmd

import (
	"fmt"
	"os"

	"github.com/7phs/coding-challenge-auction/api"
	"github.com/7phs/coding-challenge-auction/models"
	"github.com/spf13/cobra"
)

var (
	addUserName *string

	updateUserId   *int
	updateUserName *string

	getUserId *int
)

func init() {
	addUserName = UserAddCmd.Flags().String("name", "", "a users's name")

	updateUserId = UserUpdateCmd.Flags().Int("id", 0, "a user's id")
	updateUserName = UserUpdateCmd.Flags().String("name", "", "a user's name")

	getUserId = UserGetCmd.Flags().Int("id", 0, "a users's id")
}

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Commands processing a user",
}

var UserAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a user",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := api.NewUser().Add(*addUserName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			fmt.Printf("Success: #%d\n", key)
		}
	},
}

var UserUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a user",
	Run: func(cmd *cobra.Command, args []string) {
		err := api.NewUser().Update(models.UserKey(*updateUserId), *updateUserName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		} else {
			fmt.Printf("Success\n")
		}
	},
}

var UserGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a user and all bids",
	Run: func(cmd *cobra.Command, args []string) {
		item, err := api.NewUser().Get(models.UserKey(*getUserId))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		fmt.Printf("Id:   %d\n", item.Data.Id)
		fmt.Printf("Name: %s\n", item.Data.Name)
		if len(item.Data.Bids) == 0 {
			fmt.Println("Bids: no one")
		} else {
			fmt.Println("Bids:")
			for _, bid := range item.Data.Bids {
				fmt.Printf("    #%d '%s' %s %s\n", bid.ItemId, bid.ItemTitle, bid.Bid, bid.Updated)
			}
		}
	},
}
