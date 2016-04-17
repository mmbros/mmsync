package mmsync

import "sync"

// NewSimpleMutexInt ..
func NewSimpleMutexInt() MutexInt {

	acquire := func() *sync.Mutex {
		return &sync.Mutex{}
	}
	release := func(*sync.Mutex) {}

	return NewBaseMutexInt(acquire, release)
}
