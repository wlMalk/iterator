package iterator

// LimitFunc returns a modifier that stops the iterator when fn
// returns false or an error
func LimitFunc[T any](fn func(uint, T) (bool, error)) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count uint
		var finished bool
		var err error

		return &iterator[T]{
			next: func() bool {
				if err != nil {
					return false
				}

				if finished {
					return false
				}
				if !iter.Next() {
					finished = true
					return false
				}

				var item T
				item, err = iter.Get()
				if err != nil {
					finished = true
					return false
				}

				var ok bool
				ok, err = fn(count, item)
				if !ok || err != nil {
					finished = true
					return false
				}

				count++
				return true
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

// Limit returns a modifier that stops the iterator when it
// has progressed times equal to limit
func Limit[T any](limit uint) Modifier[T, T] {
	return LimitFunc(func(count uint, _ T) (bool, error) {
		return count < limit, nil
	})
}
