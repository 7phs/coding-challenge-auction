package models

import (
	"math"
	"math/rand"
	"sync"
	"testing"
)

func TestItem_Top(t *testing.T) {
	var (
		it = Item("test")

		shutdown = make(chan struct{})
		wait     sync.WaitGroup
	)

	for _, b := range []*bid{
		newBid().runQueue(shutdown, &wait).SetBid(12.34),
		newBid().runQueue(shutdown, &wait).SetBid(8.9),
		newBid().runQueue(shutdown, &wait).SetBid(1023.45),
		newBid().runQueue(shutdown, &wait).SetBid(56.78),
	} {
		user := User("user")

		it.Push(&bidRec{
			id:  newBidId(it.Id(), user.Id()),
			bid: *b,
		})
	}

	close(shutdown)
	wait.Wait()

	top, err := it.Top()
	if err != nil {
		t.Error("failed to get a top bid for an item: ", err)
	} else {
		expected := 1023.45
		exist := top.Bid()
		if math.Abs(exist-expected) > 1e-6 {
			t.Error("failed to get top bid for item. Got ", exist, ", but expected is ", expected)
		}
	}
}

func BenchmarkItem_Top(b *testing.B) {
	var (
		it = Item("test")

		shutdown = make(chan struct{})
		wait     sync.WaitGroup
		user     = User("user")
		lenBids  = 200
		bidsList = make([]*bidRec, 0, lenBids)
	)

	for i := 0; i < lenBids; i++ {
		bidsList = append(bidsList, &bidRec{
			id:  newBidId(it.Id(), user.Id()),
			bid: *newBid().runQueue(shutdown, &wait).SetBid(float64(rand.Int63n(100000000)) / precision),
		})
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		it.Push(bidsList[n%lenBids])
	}

	close(shutdown)
	wait.Wait()

	b.StopTimer()
}
