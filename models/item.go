package models

import (
	"sync"
	"math/rand"
	"container/heap"
)

const (
	defaultCapacity = 256
)

type ItemI interface {
	Id() itemKey
	SetTitle(title string)
	Title() string
	Bids(func())
}

type itemKey int32

func nextItemKey() itemKey {
	return itemKey(rand.Int31())
}

type itemInfo struct {
	sync.RWMutex

	id    itemKey
	title string
}

func (o *itemInfo) Id() itemKey {
	o.RLock()
	defer o.Unlock()

	return o.id
}

func (o *itemInfo) SetTitle(title string) {
	o.Lock()
	defer o.Unlock()

	o.title = title
}

func (o *itemInfo) Title() string {
	o.RLock()
	defer o.Unlock()

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

func (o *item) Top() BidRecI {
	o.RLock()
	defer o.RUnlock()

	return o.bidsTop[0]
}

func (o *item) Bids() []BidRecI {
	o.RLock()
	defer o.RUnlock()

	result := make([]BidRecI, 0, len(o.bidsTop))
	copy(result, o.bidsTop)

	return result
}