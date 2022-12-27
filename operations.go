package iterator

// Pipe applies the modifiers to the given iterator and returns the resulting
// iterator
func Pipe[T any](iter Iterator[T], mods ...Modifier[T, T]) Iterator[T] {
	for i := range mods {
		iter = mods[i](iter)
	}
	return iter
}

// Len exhausts the iterator to return its length
func Len[T any](iter Iterator[T]) (uint, error) {
	var count uint
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
