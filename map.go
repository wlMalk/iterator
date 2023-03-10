package iterator

import (
	"fmt"
)

// Map returns a modifier that applies fn on each item from the iterator
// It will replace the item with the value returned from fn
func Map[T any, S any](fn func(int, T) (S, error)) Modifier[T, S] {
	return FilterMap(func(i int, item T) (S, bool, error) {
		nItem, err := fn(i, item)
		if err != nil {
			return *new(S), false, err
		}
		return nItem, true, nil
	})
}

// Replace returns a modifier that changes occurances of old with new
// for as many times
func Replace[T comparable](old T, new T, times int) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count int
		return Map(func(_ int, item T) (T, error) {
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
	return Map(func(_ int, item T) (T, error) {
		if item == old {
			return new, nil
		}
		return item, nil
	})
}

// Strings is a modifier that converts all items to strings
func Strings[T any](iter Iterator[T]) Iterator[string] {
	return Map(func(_ int, item T) (string, error) {
		return fmt.Sprint(item), nil
	})(iter)
}

func ToSlices[T any](iter Iterator[Iterator[T]]) Iterator[[]T] {
	return Map(func(_ int, it Iterator[T]) ([]T, error) {
		return ToSlice(it)
	})(iter)
}
