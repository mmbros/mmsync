package mmsync

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// MAX_SLEEP_MSEC ...
const maxSleepMSec = 20

func TestSimpleMutexInt(t *testing.T) {
	const (
		totResources = 10
		totWorkers   = 1000
	)

	var (
		resources []int
		mux       MutexInt
		wg        sync.WaitGroup
	)

	mux = NewSimpleMutexInt()
	resources = make([]int, totResources)

	for j := 0; j < totWorkers; j++ {
		wg.Add(1)

		go func(wid int) {

			rid := rand.Intn(totResources)

			mux.Lock(rid)
			defer mux.Unlock(rid)

			before := resources[rid]
			if before != 0 {
				t.Errorf("worker #%03d - resources #%02d: expected %d, got %d",
					wid, rid, 0, before)
			}
			resources[rid] = 1
			time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
			after := resources[rid]
			if after != 1 {
				t.Errorf("worker #%03d - resources #%02d: expected %d, got %d",
					wid, rid, 1, after)
			}
			resources[rid] = 0

			wg.Done()
		}(j)
	}
	wg.Wait()
}

func BenchmarkSimpleMutexInt10Res(b *testing.B) {
	myBenchmarkSimpleMutexInt(b, 10)
}
func BenchmarkSimpleMutexInt50Res(b *testing.B) {
	myBenchmarkSimpleMutexInt(b, 50)
}
func BenchmarkSimpleMutexInt100Res(b *testing.B) {
	myBenchmarkSimpleMutexInt(b, 100)
}

func myBenchmarkSimpleMutexInt(b *testing.B, totResources int) {
	var (
		totWorkers = b.N
		resources  []int
		mux        MutexInt
		wg         sync.WaitGroup
	)

	mux = NewSimpleMutexInt()
	resources = make([]int, totResources)

	for j := 0; j < totWorkers; j++ {
		wg.Add(1)

		go func(wid int) {

			rid := rand.Intn(totResources)

			mux.Lock(rid)
			defer mux.Unlock(rid)

			before := resources[rid]
			if before != 0 {
				b.Errorf("worker #%03d - resources #%02d: expected %d, got %d",
					wid, rid, 0, before)
			}
			resources[rid] = 1
			time.Sleep(time.Duration(rand.Intn(maxSleepMSec)) * time.Millisecond)
			after := resources[rid]
			if after != 1 {
				b.Errorf("worker #%03d - resources #%02d: expected %d, got %d",
					wid, rid, 1, after)
			}
			resources[rid] = 0

			wg.Done()
		}(j)
	}
	wg.Wait()
}
