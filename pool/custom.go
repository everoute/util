package pool

type CustomPool[T any] struct {
	get func() T
	put func(T)
}

func NewCustomPool[T any](get func() T, put func(T)) *CustomPool[T] {
	return &CustomPool[T]{
		get: get,
		put: put,
	}
}

func (p *CustomPool[T]) Get() T {
	return p.get()
}

func (p *CustomPool[T]) Put(t T) {
	p.put(t)
}
