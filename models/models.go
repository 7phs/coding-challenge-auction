package models

import (
	log "github.com/sirupsen/logrus"
)

var (
	Users *UserTable
	Items *ItemTable
	Bids  *BidTable
)

func Init() {
	log.Info("models: init")

	Users = NewUserTable()
	Items = NewItemTable()
	Bids = NewBidTable().
		LinkStorage(Items).
		LinkStorage(Users)
}

func Shutdown() {
	log.Info("models: shutdown - start")
	Bids.Shutdown()
	log.Info("models: shutdown - finish")
}

// Using for testing
type modelsState struct {
	Users *UserTable
	Items *ItemTable
	Bids  *BidTable
}

func GetState() *modelsState {
	return &modelsState{
		Users: Users,
		Items: Items,
		Bids:  Bids,
	}
}

func PutState(state *modelsState) {
	Users = state.Users
	Items = state.Items
	Bids = state.Bids
}
