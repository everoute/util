// Why does we buffer the pool?
// 1. In sync.Pool, all objects are storaged in a linked list,
// garbage collection will scan each object in the linked list.
// 2. In sync.Pool, the New function returns a any type,
// The Pointer will be allocated on the heap.
// If the pool is large, the garbage collection will take a long time.
package pool

import (
	"sync"
)

type BufferedPool[T any] struct {
	pool   Pool[T]
	buffer []*T // a stack of objects
	mu     sync.Mutex
}

func NewBufferedPool[T any](pool Pool[T], bufferSize int) *BufferedPool[T] {
	return &BufferedPool[T]{
		pool:   pool,
		buffer: make([]*T, bufferSize),
	}
}

func (p *BufferedPool[T]) Get() *T {
	p.mu.Lock()
	if len(p.buffer) > 0 {
		t := p.buffer[len(p.buffer)-1]
		p.buffer[len(p.buffer)-1] = nil // avoid garbage collection scan
		p.buffer = p.buffer[:len(p.buffer)-1]
		p.mu.Unlock()
		return t
	}
	p.mu.Unlock()
	return p.pool.Get()
}

func (p *BufferedPool[T]) Put(t *T) {
	p.mu.Lock()
	p.buffer = append(p.buffer, t)
	p.mu.Unlock()
}
