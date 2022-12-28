package iterator

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Map returns a modifier that applies fn on each item from the iterator
// It will replace the item with the value returned from fn
func Map[T any, S any](fn func(uint, T) (S, error)) Modifier[T, S] {
	return func(iter Iterator[T]) Iterator[S] {
		var count uint
		var item S
		var err error
		return &iterator[S]{
			next: func() bool {
				if iter.Next() {
					var tItem T
					tItem, err = iter.Get()
					if err != nil {
						return false
					}
					item, err = fn(count, tItem)
					if err != nil {
						return false
					}
					count++
					return true
				}
				return false
			},
			get: func() (S, error) {
				return item, err
			},
			close: iter.Close,
			err: func() error {
				if err != nil {
					return err
				}
				return iter.Err()
			},
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

// Strings is a modifier that converts all items to strings
func Strings[T any](iter Iterator[T]) Iterator[string] {
	return Map(func(_ uint, item T) (string, error) {
		return fmt.Sprint(item), nil
	})(iter)
}

// Clamp returns a modifier to clamps items within min and max inclusively
func Clamp[T constraints.Ordered](min T, max T) Modifier[T, T] {
	return Map(func(_ uint, item T) (T, error) {
		if item < min {
			return min, nil
		} else if item > max {
			return max, nil
		}
		return item, nil
	})
}
