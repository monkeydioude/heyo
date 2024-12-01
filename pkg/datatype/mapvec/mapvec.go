package mapvec

import "github.com/monkeydioude/heyo/pkg/datatype/vec"

type MapVec[K comparable, V any] map[K]vec.Vec[V]

func New[K comparable, V any]() MapVec[K, V] {
	return make(MapVec[K, V])
}

func (mv MapVec[K, V]) Add(key K, value V) MapVec[K, V] {
	innerValue, ok := mv[key]
	if !ok {
		innerValue = vec.New[V]()
	}
	mv[key] = innerValue.Push(value)
	return mv
}
