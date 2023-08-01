package expected

import "errors"

type state int

const (
	hasNothing state = iota // default
	hasValue
	hasError
)

var (
	errNilPointer error = errors.New("nil pointer")
)

// An Expected object refers to a valid value or an error
type Expected[T any] struct {
	value T
	err   error
	state state
}

// Do not create an empty Expected object unless it will never be used.
func Empty[T any]() Expected[T] {
	return Expected[T]{
		err:   nil,
		state: hasNothing,
	}
}

// Create an Expected object from return value.
func Pack[T any](value T, err error) Expected[T] {
	if err != nil {
		return Expected[T]{
			err:   err,
			state: hasError,
		}
	}
	return Expected[T]{
		value: value,
		state: hasValue,
	}
}

// TODO: Add support for interface{...}
// Create an Expected object from return value and the value is not a nil pointer.
func PackNonNilP[P *B, B any](value P, err error) Expected[P] {
	if err != nil {
		return Expected[P]{
			err:   err,
			state: hasError,
		}
	}
	if value == nil {
		return Expected[P]{
			err:   errNilPointer,
			state: hasError,
		}
	}
	return Expected[P]{
		value: value,
		state: hasValue,
	}
}

func Value[T any](value T) Expected[T] {
	return Expected[T]{
		value: value,
		state: hasValue,
	}
}

// TODO: Add support for interface{...}
func NonNilP[T any](value *T) Expected[*T] {
	if value == nil {
		return Expected[*T]{
			err:   errNilPointer,
			state: hasError,
		}
	}
	return Expected[*T]{
		value: value,
		state: hasValue,
	}
}

// An undefined behavior will occur if you passed a nil pointer.
func Error[T any](err error) Expected[T] {
	return Expected[T]{
		err:   err,
		state: hasError,
	}
}

func Wrap[Arg, Res any](fn func(Arg) (Res, error)) func(Arg) Expected[Res] {
	return func(arg Arg) Expected[Res] {
		return Pack(fn(arg))
	}
}

// TODO: Add support for interface{...}
func WrapNonNilP[Arg any, Res *B, B any](fn func(Arg) (Res, error)) func(Arg) Expected[Res] {
	return func(arg Arg) Expected[Res] {
		return PackNonNilP(fn(arg))
	}
}

// Expect genericity begin go 1.18
func Unwrap[Arg, Res any](fn func(Arg) Expected[Res]) func(Arg) (Res, error) {
	return func(arg Arg) (Res, error) {
		return fn(arg).Get()
	}
}

func (e *Expected[T]) SetValue(value T) {
	e.value = value
	e.err = nil
	e.state = hasValue
}

func (e *Expected[T]) SetError(err error) {
	e.err = err
	e.state = hasError
}

func (e Expected[T]) Value() (T, bool) {
	return e.value, e.state == hasValue
}

func (e Expected[T]) Error() error {
	if e.state == hasError {
		return e.err
	}
	return nil
}

// TODO: release parameter type restrictions
func (e Expected[T]) ValueOr(oth T) T {
	if e.state == hasValue {
		return e.value
	}
	return oth
}

func (e Expected[T]) Get() (T, error) {
	return e.value, e.err
}

func (e Expected[T]) IsOk() bool {
	return e.state == hasValue
}

func (e Expected[T]) IsBad() bool {
	return e.state == hasError
}

// TODO: release parameter type restrictions
// call fn if it has a value
func (e Expected[T]) AndThen(fn func(T) Expected[T]) Expected[T] {
	if e.state == hasValue {
		return fn(e.value)
	}
	return e
}

// TODO: release parameter type restrictions
// call fn if it has not a value
func (e Expected[T]) OrElse(fn func(error) Expected[T]) Expected[T] {
	if e.state == hasError {
		return fn(e.err)
	}
	return e
}

// TODO: release parameter type restrictions
func (e Expected[T]) Or(res Expected[T]) Expected[T] {
	if e.state == hasValue {
		return e
	}
	return res
}

// TODO: release parameter type restrictions
func (e Expected[T]) Transform(fn func(T) T) Expected[T] {
	if e.state == hasValue {
		return Value(fn(e.value))
	}
	return e
}

// TODO: release parameter type restrictions
func (e Expected[T]) TransformError(fn func(error) error) Expected[T] {
	if e.state == hasError {
		return Error[T](fn(e.err))
	}
	return e
}
