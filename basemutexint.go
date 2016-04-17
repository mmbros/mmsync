package mmsync

import "sync"

type acquireMutexFunc func() *sync.Mutex

type releaseMutexFunc func(*sync.Mutex)

// BaseMutexInt ...
type BaseMutexInt struct {
	items        map[int]*itemType
	lock         *sync.Mutex
	acquireMutex acquireMutexFunc
	releaseMutex releaseMutexFunc
}

type itemType struct {
	lock     *sync.Mutex
	refcount int
}

// NewBaseMutexInt ...
func NewBaseMutexInt(acquire acquireMutexFunc, release releaseMutexFunc) MutexInt {
	li := &BaseMutexInt{
		items:        make(map[int]*itemType),
		lock:         &sync.Mutex{},
		acquireMutex: acquire,
		releaseMutex: release,
	}
	return li
}

// acquire ...
func (li *BaseMutexInt) acquire(id int) *sync.Mutex {
	li.lock.Lock()
	defer li.lock.Unlock()

	item, ok := li.items[id]
	if !ok {

		item = &itemType{lock: li.acquireMutex()}
		li.items[id] = item

	}
	item.refcount++
	return item.lock
}

// relaease releases the id resource.
// It is a run-time error if the id resource is not acquired.
func (li *BaseMutexInt) release(id int) *sync.Mutex {
	li.lock.Lock()
	defer li.lock.Unlock()

	item, ok := li.items[id]
	if !ok {
		panic("realese an id not acquired")
	}
	item.refcount--
	if item.refcount <= 0 {
		delete(li.items, id)
		li.releaseMutex(item.lock)
	}

	return item.lock
}

// Lock ...
func (li *BaseMutexInt) Lock(id int) {
	li.acquire(id).Lock()
}

// Unlock ...
func (li *BaseMutexInt) Unlock(id int) {
	li.release(id).Unlock()
}
