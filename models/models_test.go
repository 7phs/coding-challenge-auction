package models

import (
	"math/rand"
	"sync"
	"testing"
)

func benchModelsPush(b *testing.B, multiplier int) {
	var (
		start = make(chan struct{})
		test  sync.WaitGroup

		usersCount = 10000
		itemsCount = 200
		bidsCount  = int(30 * int64(precision))

		userTable = NewUserTable()
		itemTable = NewItemTable()

		bidTable = NewBidTable().LinkStorage(userTable, itemTable)

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

func BenchmarkModels_Push1(b *testing.B) {
	benchModelsPush(b, 1)
}

func BenchmarkModels_Push10(b *testing.B) {
	benchModelsPush(b, 10)
}

func BenchmarkModels_Push1000(b *testing.B) {
	benchModelsPush(b, 1000)
}
