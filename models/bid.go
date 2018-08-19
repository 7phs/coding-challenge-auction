package models

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	precision = float64(100000)
)

type BidRecI interface {
	Id() BidKey
	UserId() userKey
	ItemId() itemKey
	Bid() float64
	Updated() int64
}

type BidKey int64

func newBidId(userId userKey, itemId itemKey) BidKey {
	return BidKey(int64(userId)<<32 + int64(itemId))
}

func (o BidKey) UserId() userKey {
	return userKey(o >> 32)
}

func (o BidKey) ItemId() itemKey {
	return itemKey(o)
}

type bid struct {
	bid     int64
	updated int64

	queue    chan float64
	shutdown chan struct{}
}

func (o *bid) Bid() float64 {
	return float64(atomic.LoadInt64(&o.bid)) / precision
}

func (o *bid) Updated() int64 {
	return atomic.LoadInt64(&o.updated)
}

func (o *bid) SetBid(bid float64) {
	o.queue <- bid
}

func (o *bid) store(bid float64) {
	atomic.StoreInt64(&o.updated, time.Now().UnixNano())
	atomic.StoreInt64(&o.bid, int64(bid*precision))
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

func newBidRec(itemId itemKey, userId userKey, shutdown chan struct{}, wait *sync.WaitGroup) *bidRec {
	rec := &bidRec{
		id: newBidId(userId, itemId),
		bid: bid{
			queue:   make(chan float64),
			updated: time.Now().UnixNano(),
		},
	}

	rec.runQueue(shutdown, wait)

	return rec
}

func (o *bidRec) Id() BidKey {
	return o.id
}

func (o *bidRec) UserId() userKey {
	return o.id.UserId()
}

func (o *bidRec) ItemId() itemKey {
	return o.id.ItemId()
}
