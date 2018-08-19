package models

import (
	"os"
	"sync"

	"github.com/7phs/coding-challenge-auction/helpers"
)

type UserTable struct {
	users sync.Map
}

func NewUserTable() *UserTable {
	return &UserTable{}
}

func (o *UserTable) Add(user *user) error {
	_, loaded := o.users.LoadOrStore(user.Id(), user)
	if loaded {
		return os.ErrExist
	}

	return nil
}

func (o *UserTable) Get(id userKey) (UserI, error) {
	rec, ok := o.users.Load(id)
	if !ok {
		return nil, os.ErrNotExist
	}

	return rec.(UserI), nil
}

func (o *UserTable) Push(bid BidRecI) {
	id := bid.UserId()
	// try to load from cache to reduce using memory
	rec, ok := o.users.Load(id)
	if !ok {
		// it is a trade-off, just add a user's info based a bid's property
		rec, _ = o.users.LoadOrStore(id, &user{
			userInfo: userInfo{
				id: id,
			},
		})
	}

	userRec := rec.(*user)
	userRec.Push(bid)
}

// Recommended using only for testing
func (o *UserTable) Len() int {
	return helpers.SyncMapLen(&o.users)
}
