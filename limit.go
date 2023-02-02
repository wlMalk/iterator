package iterator

// Limit returns a modifier that stops the iterator when it
// has progressed times equal to limit
func Limit[T any](limit int) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count int

		return &iterator[T]{
			next: func() bool {
				if count >= limit {
					return false
				}
				if !iter.Next() {
					return false
				}

				count++
				return true
			},
			get:   iter.Get,
			close: iter.Close,
			err:   iter.Err,
		}
	}
}
