package models

import (
	"sync"

	"github.com/7phs/coding-challenge-auction/helpers"
)

type BidLinkedStorage interface {
	Push(BidRecI)
}

type BidTable struct {
	storage sync.Map

	linked []BidLinkedStorage

	shutdown chan struct{}
	wait     sync.WaitGroup
}

func NewBidTable() *BidTable {
	return &BidTable{
		shutdown: make(chan struct{}),
	}
}

// Without synchronisation, call on initial stages
func (o *BidTable) LinkStorage(linked ...BidLinkedStorage) *BidTable {
	o.linked = append(o.linked, linked...)

	return o
}

func (o *BidTable) Shutdown() {
	close(o.shutdown)
	o.wait.Wait()
}

func (o *BidTable) link(rec BidRecI) {
	for _, storage := range o.linked {
		storage.Push(rec)
	}
}

func (o *BidTable) Push(itemId ItemKey, userId UserKey, bid float64) (*bidRec, bool) {
	newRec := newBidRec(itemId, userId, o.shutdown, &o.wait)
	r, loaded := o.storage.LoadOrStore(newRec.Id(), newRec)

	rec := r.(*bidRec)
	rec.SetBid(bid)

	o.link(rec)

	return rec, !loaded
}

// Recommended using only for testing
func (o *BidTable) Len() int {
	return helpers.SyncMapLen(&o.storage)
}
