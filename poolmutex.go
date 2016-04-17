package mmsync

import (
	"errors"
	"sync"
)

// ErrClosed is the error resulting if the pool is closed via pool.Close().
var ErrClosed = errors.New("pool is closed")

// PoolMutex ...
type PoolMutex struct {
	mu    sync.Mutex
	items chan *sync.Mutex
}

// NewPoolMutex returns a new pool based on buffered channels with an initial
// capacity and maximum capacity. Factory is used when initial capacity is
// greater than zero to fill the pool. A zero initialCap doesn't fill the Pool
// until a new Get() is called. During a Get(), If there is no new connection
// available in the pool, a new connection will be created via the Factory()
// method.
func NewPoolMutex(initialCap, maxCap int) (*PoolMutex, error) {

	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	p := &PoolMutex{
		items: make(chan *sync.Mutex, maxCap),
	}

	// create initial items
	for i := 0; i < initialCap; i++ {
		p.items <- &sync.Mutex{}
	}

	return p, nil
}

func (p *PoolMutex) getItems() chan *sync.Mutex {
	p.mu.Lock()
	items := p.items
	p.mu.Unlock()
	return items
}

// Get ...
func (p *PoolMutex) Get() (*sync.Mutex, error) {
	items := p.getItems()
	if items == nil {
		return nil, ErrClosed
	}

	select {
	case item := <-items:
		if item == nil {
			return nil, ErrClosed
		}
		return item, nil
	default:
		// creates a new item
		return &sync.Mutex{}, nil
	}
}

// Put puts the item back to the pool. If the pool is full or closed,
// item is simply discarded. A nil item will be rejected.
func (p *PoolMutex) Put(item *sync.Mutex) error {
	if item == nil {
		return errors.New("item is nil. rejecting")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.items == nil {
		// pool is closed, discard passed item
		return nil
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case p.items <- item:
		return nil
	default:
		// pool is full, discard passed item
		return nil
	}
}

// Close ...
func (p *PoolMutex) Close() {
	p.mu.Lock()
	items := p.items
	p.items = nil
	p.mu.Unlock()

	if items == nil {
		return
	}

	close(items)
}

// Len ...
func (p *PoolMutex) Len() int { return len(p.getItems()) }
