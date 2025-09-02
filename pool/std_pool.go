package pool

import (
	"sync"
)

func NewStdPoll[T any](alloctor func() any) Pool[T] {
	var p stdPool[T]
	if alloctor == nil {
		p.New = func() any {
			return new(T)
		}
	} else {
		p.New = alloctor
	}
	return &p
}

type stdPool[T any] struct {
	sync.Pool
}

func (p *stdPool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *stdPool[T]) Put(t T) {
	p.Pool.Put(t)
}
