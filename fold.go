package iterator

import "golang.org/x/exp/constraints"

// Fold all items into a single value with a start value by applying fn on all items
func Fold[T any, S any](iter Iterator[T], start S, fn func(int, T, S) (S, error)) (S, error) {
	reduced := start

	_, err := Iterate(iter, func(i int, value T) (bool, error) {
		var err error
		reduced, err = fn(i, value, reduced)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return *new(S), err
	}

	return reduced, nil
}

// Reduce all items into a single value by applying fn on all items
func Reduce[T any](iter Iterator[T], fn func(int, T, T) (T, error)) (T, error) {
	if !iter.Next() {
		iter.Close()
		return *new(T), ErrNoItems
	}

	first, err := iter.Get()
	if err != nil {
		iter.Close()
		return *new(T), err
	}

	return Fold(iter, first, fn)
}

// Sum all numbers in the iterator
func Sum[T constraints.Float | constraints.Integer | constraints.Complex](iter Iterator[T]) (T, error) {
	return Fold(iter, 0, func(_ int, item T, total T) (T, error) {
		return total + item, nil
	})
}

// Concat all items with the given separator
func Concat[T ~string](iter Iterator[T], sep T) (T, error) {
	res, err := Fold(iter, sep, func(i int, item T, total T) (T, error) {
		if i == 0 {
			return total + item, nil
		}
		return total + sep + item, nil
	})
	if err != nil {
		return *new(T), err
	}

	return res[len(sep):], nil
}

// String converts and concatenates all items into a single string
func String[T any](iter Iterator[T]) (string, error) {
	return Concat(Strings(iter), "")
}
