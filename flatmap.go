package iterator

// FlatMap
func FlatMap[T any, S any](fn func(uint, uint, T) (Iterator[S], error)) Modifier[T, S] {
	return func(iter Iterator[T]) Iterator[S] {
		var finished bool
		var count uint
		var countOuter uint
		var curr Iterator[S]
		var currItem S
		var err error

		return &iterator[S]{
			next: func() bool {
				if finished || err != nil {
					return false
				}

				for curr == nil || !curr.Next() {
					if curr != nil {
						if err = curr.Err(); err != nil {
							return false
						}
						if err = curr.Close(); err != nil {
							return false
						}
					}
					if !iter.Next() {
						finished = true
						return false
					}
					var item T
					item, err = iter.Get()
					if err != nil {
						return false
					}
					curr, err = fn(countOuter, count, item)
					if err != nil {
						return false
					}
					count++
				}

				if currItem, err = curr.Get(); err != nil {
					return false
				}
				countOuter++

				return true
			},
			get: func() (S, error) {
				return currItem, err
			},
			close: func() error {
				if err := curr.Close(); err != nil {
					iter.Close()
					return err
				}
				return iter.Close()
			},
			err: func() error {
				if err != nil {
					return err
				}
				if curr != nil {
					if err := curr.Err(); err != nil {
						return err
					}
				}
				return iter.Err()
			},
		}
	}
}

// Flatten
func Flatten[T any, S Iterator[T]]() Modifier[Iterator[T], T] {
	return FlatMap(func(_ uint, _ uint, iter Iterator[T]) (Iterator[T], error) {
		return iter, nil
	})
}
