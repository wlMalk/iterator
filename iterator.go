package iterator

import (
	"errors"
)

var (
	ErrNoItems    = errors.New("iterator: no items in iterator")
)

// Iterator defines the methods needed to conform to an iterator supported by this package
type Iterator[T any] interface {
	// Next moves the iterator to the next position and reports whether
	// an item exists
	Next() bool
	// Get returns the current item
	// Multiple calls to Get without calling Next should give the same result
	Get() (T, error)
	// Close marks the iterator as done and frees resources
	Close() error
	// Err returns the first error encountered in any of the operations of the iterator
	Err() error
}

// Modifier applys an operation on the iterator and returns the resultant iterator
// The given iterator should not be used for anything else after applying the modifier
// to avoid incorrect results
type Modifier[T any, S any] func(Iterator[T]) Iterator[S]

type iterator[T any] struct {
	next  func() bool
	get   func() (T, error)
	err   func() error
	close func() error
}

func (iter *iterator[T]) Next() bool {
	if iter.next != nil {
		return iter.next()
	}
	return false
}

func (iter *iterator[T]) Get() (T, error) {
	if iter.get != nil {
		return iter.get()
	}
	return *new(T), nil
}

func (iter *iterator[T]) Close() error {
	if iter.close != nil {
		return iter.close()
	}
	return nil
}

func (iter *iterator[T]) Err() error {
	if iter.err != nil {
		return iter.err()
	}
	return nil
}
