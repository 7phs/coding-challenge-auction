package models

import (
	"testing"
	"sync"
	"math/rand"
	"sync/atomic"
)

func TestBid_Bid(t *testing.T) {
	var (
		start    = make(chan struct{})
		shutdown = make(chan struct{})
		test     sync.WaitGroup
		wait     sync.WaitGroup

		b                = newBidRec(100, 110, shutdown, &wait)
		expected int64
	)

	for i := 0; i < 10; i++ {
		test.Add(1)
		go func() {
			defer test.Done()

			var v int64

			<-start

			for i := 0; i < 1000; i++ {
				v = rand.Int63n(1000*int64(precision)) + 1
				atomic.StoreInt64(&expected, v)

				b.SetBid(float64(v) / precision)
			}
		}()
	}

	// start executing all goroutines at one time
	close(start)
	test.Wait()

	close(shutdown)
	wait.Wait()

	exist := int64(b.Bid()*precision)
	if exist!=expected {
		t.Error("failed to store the max value of bid. Got ", exist, ", but expected is ", expected)
	}
}
