package iterator

import (
	"golang.org/x/exp/constraints"
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

// Number is a constraint for all supported numeric types
type Number interface {
	constraints.Float | constraints.Integer
}
