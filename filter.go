package iterator

// Filter returns a modifier that constantly progresses the iterator to the next item
// matching pred
func Filter[T any](pred func(uint, T) (bool, error)) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count uint
		var err error

		return &iterator[T]{
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

					matches, err = pred(count, value)
					if err != nil {
						return false
					}
				}

				return matches
			},
			get:   iter.Get,
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

// RemoveFunc returns a modifier that filters away items matching fn
func RemoveFunc[T any](fn func(uint, T) (bool, error)) Modifier[T, T] {
	return Filter(func(i uint, item T) (bool, error) {
		rem, err := fn(i, item)
		if err != nil {
			return false, err
		}
		return !rem, nil
	})
}

// Remove returns a modifier that filters away items equal to rem
func Remove[T comparable](rem T) Modifier[T, T] {
	return RemoveFunc(func(_ uint, item T) (bool, error) {
		if item == rem {
			return true, nil
		}
		return false, nil
	})
}
