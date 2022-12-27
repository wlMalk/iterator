package iterator

import (
	"fmt"
)

// Map returns a modifier that applies fn on each item from the iterator
// It will replace the item with the value returned from fn
func Map[T any, S any](fn func(uint, T) (S, error)) Modifier[T, S] {
	return func(iter Iterator[T]) Iterator[S] {
		var count uint
		return &iterator[S]{
			next: func() bool {
				if iter.Next() {
					count++
					return true
				}
				return false
			},
			get: func() (S, error) {
				value, err := iter.Get()
				if err != nil {
					return *new(S), err
				}
				return fn(count, value)
			},
			close: iter.Close,
			err:   iter.Err,
		}
	}
}

// Replace returns a modifier that changes occurances of old with new
// for as many times
func Replace[T comparable](old T, new T, times uint) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count uint
		return Map(func(_ uint, item T) (T, error) {
			if count < times && item == old {
				count++
				return new, nil
			}
			return item, nil
		})(iter)
	}
}

// ReplaceAll returns a modifier that changes all occurances of old with new
func ReplaceAll[T comparable](old T, new T) Modifier[T, T] {
	return Map(func(_ uint, item T) (T, error) {
		if item == old {
			return new, nil
		}
		return item, nil
	})
}

// Strings returns a modifier that converts all items to strings
func Strings[T any](iter Iterator[T]) Iterator[string] {
	return Map(func(_ uint, item T) (string, error) {
		return fmt.Sprint(item), nil
	})(iter)
}
