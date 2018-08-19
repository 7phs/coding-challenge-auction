package models

import (
	"sync"
	"os"
)

type UserTable struct {
	users sync.Map
}

func NewUserTable() *UserTable {
	return &UserTable{}
}

func (o *UserTable) Add(item *item) error {
	_, loaded := o.users.LoadOrStore(item.Id(), item)
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

	itemRec := rec.(*item)
	itemRec.Push(bid)
}
