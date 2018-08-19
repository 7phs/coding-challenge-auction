package models

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/7phs/coding-challenge-auction/helpers"
)

type testLinkedStorage struct {
	items sync.Map
	users sync.Map
}

func newTestLinkedStorage() *testLinkedStorage {
	return &testLinkedStorage{}
}

func (o *testLinkedStorage) Push(b BidRecI) {
	o.items.Store(b.ItemId(), b)
	o.users.Store(b.UserId(), b)
}

func (o *testLinkedStorage) ItemsLen() int {
	return helpers.SyncMapLen(&o.items)
}

func (o *testLinkedStorage) UsersLen() int {
	return helpers.SyncMapLen(&o.users)
}

func TestBidTable_Push(t *testing.T) {
	var (
		start = make(chan struct{})
		test  sync.WaitGroup

		usersCount = 200
		itemsCount = 10
		bidsCount  = int(30 * int64(precision))

		testLinked = newTestLinkedStorage()

		bidTable = NewBidTable().
				LinkStorage(testLinked)
	)

	for i := 0; i < 5; i++ {
		test.Add(1)
		go func() {
			defer test.Done()

			users := rand.Perm(usersCount)
			items := rand.Perm(itemsCount)
			bids := rand.Perm(bidsCount)

			<-start

			for j := 0; j < usersCount; j++ {
				for k := 0; k < itemsCount; k++ {
					bidTable.Push(itemKey(items[k]+1), userKey(users[j]+1), float64(bids[(i+k)%bidsCount])/precision)
				}
			}
		}()
	}

	close(start)
	test.Wait()

	bidTable.Shutdown()

	if existLen := testLinked.ItemsLen(); existLen != itemsCount {
		t.Error("failed to store items. Got ", existLen, ", but expected is ", itemsCount)
	}

	if existLen := testLinked.UsersLen(); existLen != usersCount {
		t.Error("failed to store users. Got ", existLen, ", but expected is ", usersCount)
	}

	expectedLen := itemsCount * usersCount
	if existLen := bidTable.Len(); existLen != expectedLen {
		t.Error("failed to store bids. Got ", existLen, ", but expected is ", expectedLen)
	}
}

func benchBidTablePush(b *testing.B, multiplier int) {
	var (
		start = make(chan struct{})
		test  sync.WaitGroup

		usersCount = 1000
		itemsCount = 20
		bidsCount  = int(30 * int64(precision))

		testLinked = newTestLinkedStorage()

		bidTable = NewBidTable().
				LinkStorage(testLinked)

		users = rand.Perm(usersCount)
		items = rand.Perm(itemsCount)
		bids  = rand.Perm(bidsCount)
	)

	for i := 0; i < multiplier; i++ {
		test.Add(1)
		go func() {
			defer test.Done()

			<-start

			for n := 1; ; n++ {
				bidTable.Push(itemKey(items[n%itemsCount]), userKey(users[n%usersCount]), float64(bids[n%bidsCount])/precision)

				if n > b.N/multiplier {
					break
				}
			}
		}()
	}

	b.ResetTimer()

	close(start)
	test.Wait()

	b.StopTimer()

	bidTable.Shutdown()
}

func BenchmarkBidTable_Push1(b *testing.B) {
	benchBidTablePush(b, 1)
}

func BenchmarkBidTable_Push10(b *testing.B) {
	benchBidTablePush(b, 10)
}

func BenchmarkBidTable_Push1000(b *testing.B) {
	benchBidTablePush(b, 1000)
}
