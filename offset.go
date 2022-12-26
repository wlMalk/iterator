package iterator

// OffsetFunc returns a modifier that progresses the iterator as long as fn
// returns true
func OffsetFunc[T any](fn func(uint64, T) (bool, error)) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count uint64
		var started bool
		var err error

		return &iterator[T]{
			next: func() bool {
				if err != nil {
					return false
				}

				if started {
					return iter.Next()
				} else {

					for {
						if !iter.Next() {
							return false
						}

						var item T
						item, err = iter.Get()
						if err != nil {
							return false
						}

						var skip bool
						skip, err = fn(count, item)
						if err != nil {
							return false
						}
						if !skip {
							started = true
							return true
						}

						count++
					}
				}
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

// Offset returns a modifier that progresses the iterator times equal
// to offset
func Offset[T any](offset uint64) Modifier[T, T] {
	return OffsetFunc(func(count uint64, _ T) (bool, error) {
		return count < offset, nil
	})
}
