package models

import (
	"container/heap"
	"math/rand"
	"os"
	"sync"
)

const (
	defaultCapacity = 256
)

type ItemI interface {
	Id() ItemKey
	SetTitle(title string)
	Title() string
	Top() (BidRecI, error)
	Bids() []BidRecI
}

type ItemKey int32

func nextItemKey() ItemKey {
	return ItemKey(rand.Int31())
}

type itemInfo struct {
	sync.RWMutex

	id    ItemKey
	title string
}

func (o *itemInfo) Id() ItemKey {
	o.RLock()
	defer o.RUnlock()

	return o.id
}

func (o *itemInfo) SetTitle(title string) {
	o.Lock()
	defer o.Unlock()

	o.title = title
}

func (o *itemInfo) Title() string {
	o.RLock()
	defer o.RUnlock()

	return o.title
}

type item struct {
	sync.RWMutex
	itemInfo

	bids    sync.Map
	bidsTop bidHeap
}

func Item(title string) *item {
	return &item{
		itemInfo: itemInfo{
			id:    nextItemKey(),
			title: title,
		},
		bidsTop: make(bidHeap, 0, defaultCapacity),
	}
}

func (o *item) Push(bid BidRecI) {
	_, loaded := o.bids.LoadOrStore(bid.Id(), bid)

	if !loaded {
		o.Lock()
		heap.Push(&o.bidsTop, bid)
		o.Unlock()
	}
}

func (o *item) Top() (BidRecI, error) {
	o.RLock()
	defer o.RUnlock()

	if o.bidsTop.Len() == 0 {
		return nil, os.ErrNotExist
	}

	return o.bidsTop[0], nil
}

func (o *item) Bids() []BidRecI {
	o.RLock()
	defer o.RUnlock()

	result := make([]BidRecI, 0, len(o.bidsTop))
	copy(result, o.bidsTop)

	return result
}
