package iterator

// FilterMap returns a modifier that constantly progresses the iterator to the next item
// matching fn, and transforms it into an iterator for a different type.
func FilterMap[T any, S any](fn func(int, T) (S, bool, error)) Modifier[T, S] {
	return filterMap(func(index int, item T) (mappedItem S, matches bool, canContinue bool, err error) {
		canContinue = true
		mappedItem, matches, err = fn(index, item)
		return
	})
}

// TakeWhile returns a modifier which makes the iterator return items
// matching pred and stops at the first item that does not.
func TakeWhile[T any](pred func(int, T) (bool, error)) Modifier[T, T] {
	return filter(func(index int, item T) (bool, bool, error) {
		matches, err := pred(index, item)
		if !matches || err != nil {
			return false, false, err
		}

		return true, true, nil
	})
}

// Slice returns a modifier which makes the iterator return
// items in the range [from, to) increasing by step.
func Slice[T any](from, to, step int) Modifier[T, T] {
	return filter(func(index int, item T) (bool, bool, error) {
		if index < from {
			return false, true, nil
		}

		if to > -1 && index >= to {
			return false, false, nil
		}

		if step == 1 {
			// fmt.Println(index, item)
			return true, to == -1 || index+step < to, nil
		}

		return (index-from)%step == 0, to == -1 || index+step < to, nil
	})
}

// DropWhile returns a modifier which makes the iterator drop items
// matching pred and starts with the first item that does not.
func DropWhile[T any](pred func(int, T) (bool, error)) Modifier[T, T] {
	var stoppedDropping bool
	return filter(func(index int, item T) (bool, bool, error) {
		if stoppedDropping {
			return true, true, nil
		}

		dropped, err := pred(index, item)
		if err != nil {
			return false, false, err
		}
		if dropped {
			return false, true, nil
		}

		stoppedDropping = true
		return true, true, nil
	})
}

// Splice returns a modifier which replaces items in the range [from, to) with all items from injected.
func Splice[T any](from, to int, injected Iterator[T]) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		if from != to {
			iter = filter(func(index int, _ T) (bool, bool, error) {
				if index >= from && (to == -1 || index < to) {
					return false, to != -1, nil
				}
				return true, true, nil
			})(iter)
		}
		if injected != nil {
			iter = Inject(from, injected)(iter)
		}
		return iter
	}
}

func filter[T any](fn func(int, T) (matches bool, canContinue bool, err error)) Modifier[T, T] {
	return filterMap(func(index int, item T) (T, bool, bool, error) {
		matches, canContinue, err := fn(index, item)
		if err != nil {
			return *new(T), false, false, err
		}
		if !matches {
			return *new(T), false, canContinue, nil
		}
		return item, true, canContinue, nil
	})
}

func filterMap[T any, S any](fn func(int, T) (item S, matches bool, canContinue bool, err error)) Modifier[T, S] {
	return func(iter Iterator[T]) Iterator[S] {
		var count int
		var curr S
		var done bool
		var err error

		return &iterator[S]{
			next: func() bool {
				var matches bool
				var value T

				for !done && !matches {
					if !iter.Next() {
						return false
					}
					count++
					value, err = iter.Get()
					if err != nil {
						return false
					}
					var canContinue bool
					curr, matches, canContinue, err = fn(count-1, value)
					if err != nil {
						return false
					}
					if !canContinue {
						done = true
					}
				}

				return matches
			},
			get: func() (S, error) {
				if err != nil {
					return *new(S), err
				}
				return curr, nil
			},
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
