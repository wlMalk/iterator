package iterator

// ContainsFunc calls fn on each item in the given iterator.
// It stops as soon as fn returns true and reports that.
func ContainsFunc[T any](iter Iterator[T], fn func(uint, T) (bool, error)) (bool, error) {
	result := false
	if _, err := Iterate(iter, func(i uint, item T) (bool, error) {
		contains, err := fn(i, item)
		if err != nil {
			return false, err
		}
		if contains {
			result = true
		}
		return !contains, nil
	}); err != nil {
		return false, err
	}
	return result, nil
}

// ContainsAny returns true as soon as it encounters any of the given values in the given iterator.
func ContainsAny[T comparable](iter Iterator[T], values ...T) (bool, error) {
	return ContainsFunc(iter, func(_ uint, item T) (bool, error) {
		for _, value := range values {
			if item == value {
				return true, nil
			}
		}
		return false, nil
	})
}

// ContainsAll returns true when it has encountered all the given values in the given iterator.
func ContainsAll[T comparable](iter Iterator[T], values ...T) (bool, error) {
	return ContainsFunc(iter, func(_ uint, item T) (bool, error) {
		for i, value := range values {
			if item == value {
				if len(values) == 1 {
					return true, nil
				}
				values = append(values[:i], values[i+1:]...)
				return false, nil
			}
		}
		return false, nil
	})
}
