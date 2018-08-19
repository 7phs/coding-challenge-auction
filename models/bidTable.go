package models

import "sync"

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
func (o *BidTable) LinkStorage(linked ... BidLinkedStorage) *BidTable {
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

func (o *BidTable) Push(itemId itemKey, userId userKey, bid float64) *bidRec {
	newRec := newBidRec(itemId, userId, o.shutdown, &o.wait)

	r, _ := o.storage.LoadOrStore(newRec.Id(), newRec)

	rec := r.(*bidRec)
	rec.SetBid(bid)

	o.link(rec)

	return rec
}
