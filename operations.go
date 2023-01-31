package iterator

// Pipe applies the modifiers to the given iterator and returns the resulting
// iterator
func Pipe[T any](iter Iterator[T], mods ...Modifier[T, T]) Iterator[T] {
	for i := range mods {
		iter = mods[i](iter)
	}
	return iter
}

func Iterate[T any](iterator Iterator[T], f func(int, T) (bool, error)) (int, error) {
	var count int
	shouldContinue := true
	for shouldContinue && iterator.Next() {
		r, err := iterator.Get()
		if err != nil {
			return count, err
		}
		shouldContinue, err = f(count, r)
		if err != nil {
			return count, err
		}
		count++
	}

	if err := iterator.Err(); err != nil {
		return count, err
	}

	if err := iterator.Close(); err != nil {
		return count, err
	}

	return count, nil
}

// Len exhausts the iterator to return its length
func Len[T any](iter Iterator[T]) (int, error) {
	var count int
	for iter.Next() {
		count++
	}
	return count, iter.Err()
}

// Equal reports whether the given iterators are equal
func Equal[T comparable](iters ...Iterator[T]) (bool, error) {
	for {
		var firstHasNext bool
		var firstItem T

		for i := range iters {
			if i == 0 {
				firstHasNext = iters[i].Next()
				if !firstHasNext {
					if err := iters[i].Err(); err != nil {
						return false, err
					}
				}

				var err error
				firstItem, err = iters[i].Get()
				if err != nil {
					return false, err
				}
			} else {
				if iters[i].Next() != firstHasNext {
					if firstHasNext {
						return false, iters[i].Err()
					}
					return false, nil
				}

				item, err := iters[i].Get()
				if err != nil {
					return false, err
				}

				if item != firstItem {
					return false, nil
				}
			}
		}

		if !firstHasNext {
			var err error
			for i := range iters {
				if closeErr := iters[i].Close(); err == nil && closeErr != nil {
					err = closeErr
				}
			}
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
}

// One returns the only item of an iterator.
// If the iterator has no items then it returns ErrNoItems,
// and if the iterator has more than one item then it returns ErrMultiItems.
func One[T any](iterator Iterator[T]) (T, error) {
	defer iterator.Close()

	var t T
	var err error
	if iterator.Next() {
		t, err = iterator.Get()
		if err != nil {
			return *new(T), err
		}
		if iterator.Next() {
			return *new(T), ErrMultiItems
		}
	} else {
		return *new(T), ErrNoItems
	}

	if err := iterator.Err(); err != nil {
		return *new(T), err
	}

	if err := iterator.Close(); err != nil {
		return *new(T), err
	}

	return t, nil
}

// Async calls fn with the given iterator in a separate goroutine
// and returns the result in a ValErr channel.
func Async[T any, V any](iterator Iterator[T], fn func(Iterator[T]) (V, error)) <-chan ValErr[V] {
	c := make(chan ValErr[V])
	go func() {
		v, err := fn(iterator)
		c <- ValErr[V]{Val: v, Err: err}
		close(c)
	}()
	return c
}
