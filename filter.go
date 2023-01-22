package iterator

// Filter returns a modifier that constantly progresses the iterator to the next item
// matching pred
func Filter[T any](pred func(uint, T) (bool, error)) Modifier[T, T] {
	return FilterMap(func(i uint, item T) (T, bool, error) {
		matches, err := pred(i, item)
		if err != nil {
			return *new(T), false, err
		}
		if !matches {
			return *new(T), false, err
		}
		return item, matches, nil
	})
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

func DistinctFunc[T any, S comparable](fn func(uint, T) (S, error)) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		set := make(map[S]struct{})
		return RemoveFunc(func(i uint, item T) (bool, error) {
			key, err := fn(i, item)
			if err != nil {
				return false, err
			}
			_, ok := set[key]
			if ok {
				return true, nil
			}
			set[key] = struct{}{}
			return false, nil
		})(iter)
	}
}

// Distinct is a modifier that skips duplicate items
func Distinct[T comparable](iter Iterator[T]) Iterator[T] {
	return DistinctFunc(func(_ uint, item T) (T, error) {
		return item, nil
	})(iter)
}
