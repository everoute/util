package pool

import (
	"sync"
)

func NewStdPoll[T any](alloctor func() any) Pool[T] {
	var p stdPool[T]
	if alloctor == nil {
		p.pool.New = func() any {
			return new(T)
		}
	} else {
		p.pool.New = alloctor
	}
	return &p
}

type stdPool[T any] struct {
	pool sync.Pool
}

func (p *stdPool[T]) Get() *T {
	return p.pool.Get().(*T)
}

func (p *stdPool[T]) Put(t *T) {
	p.pool.Put(t)
}
