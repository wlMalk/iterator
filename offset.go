package iterator

// Offset returns a modifier that progresses the iterator times equal
// to offset
func Offset[T any](offset int) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count int

		return &iterator[T]{
			next: func() bool {
				if count > offset {
					return iter.Next()
				}
				for count <= offset {
					if !iter.Next() {
						return false
					}

					count++
				}
				return true
			},
			get:   iter.Get,
			close: iter.Close,
			err:   iter.Err,
		}
	}
}
