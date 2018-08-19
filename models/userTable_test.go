package models

import (
	"math/rand"
	"sync"
	"testing"
)

func TestUserTable_Push(t *testing.T) {
	var (
		start    = make(chan struct{})
		shutdown = make(chan struct{})
		test     sync.WaitGroup
		wait     sync.WaitGroup

		multiplier = 5
		itemsCount = 200
		bidsCount  = int(30 * int64(Precision))

		userTable = NewUserTable()

		items = rand.Perm(itemsCount)
		bids  = rand.Perm(bidsCount)
	)

	for i := 0; i < multiplier; i++ {
		test.Add(1)
		go func() {
			defer test.Done()

			<-start

			for j := 0; j < itemsCount; j++ {
				user := User("")
				if err := userTable.Add(user); err != nil {
					t.Error("failed to add a user: ", err)
				}

				b := newBidRec(ItemKey(items[j]), user.Id(), shutdown, &wait)
				b.SetBid(float64(bids[j%bidsCount]) / Precision)

				userTable.Push(b)
			}
		}()
	}

	close(start)
	test.Wait()

	expectedLen := multiplier * itemsCount
	if existLen := userTable.Len(); existLen != expectedLen {
		t.Error("failed to store users. Got ", existLen, ", but expected is ", expectedLen)
	}
}
