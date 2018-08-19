package models

import (
	"os"
	"sync"

	"github.com/7phs/coding-challenge-auction/helpers"
)

type ItemTable struct {
	items sync.Map
}

func NewItemTable() *ItemTable {
	return &ItemTable{}
}

func (o *ItemTable) Add(item *item) error {
	_, loaded := o.items.LoadOrStore(item.Id(), item)
	if loaded {
		return os.ErrExist
	}

	return nil
}

func (o *ItemTable) Get(id itemKey) (ItemI, error) {
	rec, ok := o.items.Load(id)
	if !ok {
		return nil, os.ErrNotExist
	}

	return rec.(ItemI), nil
}

func (o *ItemTable) Push(bid BidRecI) {
	key := bid.ItemId()
	// try to load from cache to reduce using memory
	rec, ok := o.items.Load(key)
	if !ok {
		// it is a trade-off, just add a item's info based a bid's property
		rec, _ = o.items.LoadOrStore(key, &item{
			itemInfo: itemInfo{
				id: key,
			},
		})
	}

	itemRec := rec.(*item)
	itemRec.Push(bid)
}

// Recommended using only for testing
func (o *ItemTable) Len() int {
	return helpers.SyncMapLen(&o.items)
}
