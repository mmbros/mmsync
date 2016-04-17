package mmsync

import "sync"

// NewPoolMutexInt ...
func NewPoolMutexInt(pool *PoolMutex) MutexInt {
	if pool == nil {
		panic("PoolMutex is nil")
	}

	acquire := func() *sync.Mutex {
		mu, err := pool.Get()
		if err != nil {
			panic("can't get a mutex from pool")
		}
		return mu
	}
	release := func(mu *sync.Mutex) {
		pool.Put(mu)
	}

	return NewBaseMutexInt(acquire, release)
}
