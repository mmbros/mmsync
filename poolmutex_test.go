package mmsync

import (
	"sync"
	"testing"
)

var (
	InitialCap = 10
	MaximumCap = 20
)

func TestNewPoolMutex(t *testing.T) {
	var err error
	_, err = NewPoolMutex(InitialCap, MaximumCap)
	if err != nil {
		t.Errorf("NewPoolMutex error: %s", err)
	}

	_, err = NewPoolMutex(20, 10)
	if err == nil {
		t.Errorf("NewPoolMutex error: %s", "no error when initialCap > maxCap")
	}

	_, err = NewPoolMutex(-1, MaximumCap)
	if err == nil {
		t.Errorf("NewPoolMutex error: %s", "no error when initialCap < 0")
	}

	_, err = NewPoolMutex(0, 0)
	if err == nil {
		t.Errorf("NewPoolMutex error: %s", "no error when maxCap <= 0")
	}
}

func TestPoolMutex_Get(t *testing.T) {
	p, _ := NewPoolMutex(InitialCap, MaximumCap)
	defer p.Close()

	_, err := p.Get()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	// after one get, current capacity should be lowered by one.
	if p.Len() != (InitialCap - 1) {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.Len())
	}

	// get them all
	var wg sync.WaitGroup
	for i := 0; i < (InitialCap - 1); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Get()
			if err != nil {
				t.Errorf("Get error: %s", err)
			}
		}()
	}
	wg.Wait()

	if p.Len() != 0 {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.Len())
	}

	_, err = p.Get()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
}

func TestPoolMutex_Put(t *testing.T) {
	p, err := NewPoolMutex(InitialCap, MaximumCap)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	// get/create from the pool
	items := make([]*sync.Mutex, MaximumCap)
	for i := 0; i < MaximumCap; i++ {
		item, _ := p.Get()
		items[i] = item
	}

	// now put them all back
	for _, item := range items {
		p.Put(item)
	}

	if p.Len() != MaximumCap {
		t.Errorf("Put error len. Expecting %d, got %d",
			1, p.Len())
	}

	item, _ := p.Get()
	p.Close() // close pool

	p.Put(item) // try to put into a full pool
	if p.Len() != 0 {
		t.Errorf("Put error. Closed pool shouldn't allow to put items.")
	}
}
