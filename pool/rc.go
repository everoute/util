package pool

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/everoute/util/nocopy"
)

// RCPool is a simple implementation of reference counting pool.
// It is not implemented the Pool interface, it is a wrapper of the Pool interface only.
type RCPool[T any] struct {
	Pool[T]
}

// NewRCPool creates a new RCPool.
func NewRCPool[T any](pool Pool[T]) *RCPool[T] {
	return &RCPool[T]{
		Pool: pool,
	}
}

// rcPoolCell is the cell of the RC object. It is used to store the reference count and the value.
type rcPoolCell[T any] struct {
	count int64
	pool  Pool[T]
	value T
}

// RC is the reference counting object.
// Do not clone this object!!!
type RC[T any] struct {
	nocopy.NoCopy
	*rcPoolCell[T]
}

func (p *RCPool[T]) New() RC[T] {
	return RC[T]{
		rcPoolCell: &rcPoolCell[T]{
			count: 1,
			pool:  p.Pool,
			value: p.Get(),
		},
	}
}

func (rc *RC[T]) Ref() RC[T] {
	_ = atomic.AddInt64(&rc.count, 1)
	return RC[T]{
		rcPoolCell: rc.rcPoolCell,
	}
}

// Do not use the rc object after unref!
func (rc *RC[T]) Unref() {
	nc := atomic.AddInt64(&rc.count, -1) // new count
	if nc == 0 {
		rc.pool.Put(rc.value)
	}
	if nc < 0 {
		panic("RC: reference count is less than 0")
	}
}

func init() {
	// Assert object size
	sizeOfRC := unsafe.Sizeof(RC[any]{})
	sizeOfPointer := unsafe.Sizeof((*int)(nil))
	if sizeOfRC != sizeOfPointer {
		panic(fmt.Sprintf("RC: size of RC is not a pointer, sizeOfRC: %d, sizeOfPointer: %d", sizeOfRC, sizeOfPointer))
	}
}
