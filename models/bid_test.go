package models

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

func TestBid_Bid(t *testing.T) {
	var (
		start    = make(chan struct{})
		shutdown = make(chan struct{})
		test     sync.WaitGroup
		wait     sync.WaitGroup

		b        = newBidRec(100, 110, shutdown, &wait)
		expected int64
	)

	for i := 0; i < 10; i++ {
		test.Add(1)
		go func() {
			defer test.Done()

			var v int64

			<-start

			for i := 0; i < 1000; i++ {
				v = rand.Int63n(1000*int64(Precision)) + 1
				atomic.StoreInt64(&expected, v)

				b.SetBid(float64(v) / Precision)
			}
		}()
	}

	// start executing all goroutines at one time
	close(start)
	test.Wait()

	close(shutdown)
	wait.Wait()

	exist := int64(b.Bid() * Precision)
	if exist != expected {
		t.Error("failed to store the max value of bid. Got ", exist, ", but expected is ", expected)
	}
}

func benchBidBid(b *testing.B, multiplier int) {
	var (
		start    = make(chan struct{})
		shutdown = make(chan struct{})
		wait     sync.WaitGroup
		test     sync.WaitGroup

		bi = newBidRec(100, 110, shutdown, &wait)
	)

	for i := 0; i < multiplier; i++ {
		test.Add(1)
		go func() {
			defer test.Done()

			v := rand.Int63n(1000*int64(Precision)) + 1

			<-start

			for n := 0; n < b.N/multiplier; n++ {
				bi.SetBid(float64(v) / Precision)
			}
		}()
	}

	b.ResetTimer()

	close(start)
	test.Wait()
	close(shutdown)
	wait.Wait()

	b.StopTimer()
}

func BenchmarkBid_Bid1(b *testing.B) {
	benchBidBid(b, 1)
}

func BenchmarkBid_Bid10(b *testing.B) {
	benchBidBid(b, 10)
}

func BenchmarkBid_Bid1000(b *testing.B) {
	benchBidBid(b, 1000)
}
