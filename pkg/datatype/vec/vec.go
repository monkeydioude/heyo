package vec

import "slices"

type Vec[T any] []T

func New[T any]() Vec[T] {
	return make(Vec[T], 0)
}

func (Vec Vec[T]) Push(value T) Vec[T] {
	return append(Vec, value)
}

func (Vec Vec[T]) DeleteFunc(value T, condFn func(T) bool) Vec[T] {
	return slices.DeleteFunc(Vec, condFn)
}
