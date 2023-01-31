package iterator

// FilterMap returns a modifier that constantly progresses the iterator to the next item
// matching fn, and transforms it into an iterator for a different type.
func FilterMap[T any, S any](fn func(int, T) (S, bool, error)) Modifier[T, S] {
	return func(iter Iterator[T]) Iterator[S] {
		var count int
		var curr S
		var err error

		return &iterator[S]{
			next: func() bool {
				var matches bool
				var value T

				for !matches {
					if !iter.Next() {
						return false
					}
					count++
					value, err = iter.Get()
					if err != nil {
						return false
					}

					curr, matches, err = fn(count, value)
					if err != nil {
						return false
					}
				}

				return matches
			},
			get: func() (S, error) {
				if err != nil {
					return *new(S), err
				}
				return curr, nil
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
