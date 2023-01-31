package iterator

// FlatMap returns a modifier to map items to iterators and flattens them into a single iterator
func FlatMap[T any, S any](fn func(int, int, T) (Iterator[S], error)) Modifier[T, S] {
	return func(iter Iterator[T]) Iterator[S] {
		var finished bool
		var count int
		var countOuter int
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

// Flatten is a modifier that applies on iterator of iterators
// It flattens them into a single iterator
func Flatten[T any](iter Iterator[Iterator[T]]) Iterator[T] {
	return FlatMap(func(_ int, _ int, iter Iterator[T]) (Iterator[T], error) {
		return iter, nil
	})(iter)
}
