package models

var (
	Users *UserTable
	Items *ItemTable
	Bids  *BidTable
)

func Init() {
	Users = NewUserTable()
	Items = NewItemTable()
	Bids = NewBidTable().
		LinkStorage(Items).
		LinkStorage(Users)
}

func Shutdown() {
	Bids.Shutdown()
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
