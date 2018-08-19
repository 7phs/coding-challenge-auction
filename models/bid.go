package models

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Precision = float64(100000)
)

type BidRecI interface {
	Id() BidKey
	UserId() UserKey
	ItemId() ItemKey
	Bid() float64
	Updated() int64
}

type BidKey int64

func newBidId(itemId ItemKey, userId UserKey) BidKey {
	return BidKey(int64(itemId)<<32 + int64(userId))
}

func (o BidKey) UserId() UserKey {
	return UserKey(o)
}

func (o BidKey) ItemId() ItemKey {
	return ItemKey(o >> 32)
}

type bid struct {
	bid     int64
	updated int64

	queue    chan float64
	shutdown chan struct{}
}

func newBid() *bid {
	return &bid{
		queue: make(chan float64),
	}
}

func (o *bid) Bid() float64 {
	return float64(atomic.LoadInt64(&o.bid)) / Precision
}

func (o *bid) Updated() int64 {
	return atomic.LoadInt64(&o.updated)
}

func (o *bid) SetBid(bid float64) *bid {
	o.queue <- bid

	return o
}

func (o *bid) store(bid float64) {
	atomic.StoreInt64(&o.updated, time.Now().UnixNano())
	atomic.StoreInt64(&o.bid, int64(bid*Precision))
}

func (o *bid) runQueue(shutdown chan struct{}, wait *sync.WaitGroup) *bid {
	o.shutdown = shutdown

	wait.Add(1)
	go func() {
		defer func() {
			wait.Done()
		}()

		for {
			select {
			case b := <-o.queue:
				o.store(b)

			case <-o.shutdown:
				return
			}
		}
	}()

	return o
}

type bidRec struct {
	bid

	id BidKey
}

func newBidRec(itemId ItemKey, userId UserKey, shutdown chan struct{}, wait *sync.WaitGroup) *bidRec {
	rec := &bidRec{
		id:  newBidId(itemId, userId),
		bid: *newBid(),
	}

	rec.runQueue(shutdown, wait)

	return rec
}

func (o *bidRec) String() string {
	return fmt.Sprintf("#%d; item: #%d; user: #%d; bid: %.5f", o.Id(), o.ItemId(), o.UserId(), o.Bid())
}

func (o *bidRec) Id() BidKey {
	return o.id
}

func (o *bidRec) UserId() UserKey {
	return o.id.UserId()
}

func (o *bidRec) ItemId() ItemKey {
	return o.id.ItemId()
}

func (o *bidRec) SetBid(bid float64) {
	o.bid.SetBid(bid)
}
