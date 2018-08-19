package models

import (
	"math/rand"
	"sync"
	"testing"
)

func TestItemTable_Push(t *testing.T) {
	var (
		start    = make(chan struct{})
		shutdown = make(chan struct{})
		test     sync.WaitGroup
		wait     sync.WaitGroup

		multiplier = 5
		usersCount = 200
		bidsCount  = int(30 * int64(precision))

		itemTable = NewItemTable()

		users = rand.Perm(usersCount)
		bids  = rand.Perm(bidsCount)
	)

	for i := 0; i < multiplier; i++ {
		test.Add(1)
		go func() {
			defer test.Done()

			<-start

			for j := 0; j < usersCount; j++ {
				item := Item("")
				if err := itemTable.Add(item); err != nil {
					t.Error("failed to add an item: ", err)
				}

				b := newBidRec(item.Id(), userKey(users[j]), shutdown, &wait)
				b.SetBid(float64(bids[j%bidsCount]) / precision)

				itemTable.Push(b)
			}
		}()
	}

	close(start)
	test.Wait()

	expectedLen := multiplier * usersCount
	if existLen := itemTable.Len(); existLen != expectedLen {
		t.Error("failed to store items. Got ", existLen, ", but expected is ", expectedLen)
	}
}
