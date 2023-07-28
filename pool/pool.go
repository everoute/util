package pool

type Putter[T any] interface {
	Put(*T)
}

type Getter[T any] interface {
	Get() *T
}

type Pool[T any] interface {
	Putter[T]
	Getter[T]
}
