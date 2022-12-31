package iterator

// Join multiple iterators to form a larger iterator
func Join[T any](iters ...Iterator[T]) Iterator[T] {
	var curr int
	var finished bool

	return &iterator[T]{
		next: func() bool {
			if finished {
				return false
			}

			if len(iters) == 0 {
				return false
			}

			for !iters[curr].Next() {
				if curr == len(iters)-1 {
					finished = true
					return false
				}

				curr++
			}

			return true
		},
		get: func() (T, error) {
			return iters[curr].Get()
		},
		close: func() error {
			var err error
			for i := range iters {
				if closeErr := iters[i].Close(); err == nil && closeErr != nil {
					err = closeErr
				}
			}
			return err
		},
		err: func() error {
			for i := range iters {
				if err := iters[i].Err(); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// JoinLeading returns a modifier to join an iterator to the end of the given leading iterator
func JoinLeading[T any](leading Iterator[T]) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		return Join(leading, iter)
	}
}

// JoinTrailing returns a modifier to join the given trailing iterator to the end of an iterator
func JoinTrailing[T any](trailing Iterator[T]) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		return Join(iter, trailing)
	}
}

// Append returns a modifier to append items to an iterator
func Append[T any](items ...T) Modifier[T, T] {
	return JoinTrailing(FromSlice(items))
}

// Prepend returns a modifier to prepend items to an iterator
func Prepend[T any](items ...T) Modifier[T, T] {
	return JoinLeading(FromSlice(items))
}
