package iterator

// SplitFunc returns a modifier that splits the iterator when fn returns true
// The return values of fn also determine whether the item should be part of the firat part, second part, both or neither
func SplitFunc[T any](fn func(int, T) (split, inA, inB bool, err error)) Modifier[T, Iterator[T]] {
	return func(iter Iterator[T]) Iterator[Iterator[T]] {
		var finished bool
		var curr Iterator[T]
		var err error
		var item T
		var startWithItem bool
		var count int

		next := func() func() (T, bool, error) {
			var iteratorFinished bool
			return func() (T, bool, error) {
				if finished || iteratorFinished || err != nil {
					return *new(T), false, err
				}

				if startWithItem {
					startWithItem = false
					return item, true, nil
				}

				if !iter.Next() {
					finished = true
					return *new(T), false, err
				}

				item, err = iter.Get()
				if err != nil {
					return *new(T), false, err
				}

				shouldSplit, inA, inB, fnErr := fn(count, item)
				if fnErr != nil {
					err = fnErr
					return *new(T), false, err
				}
				count++

				if !shouldSplit {
					return item, true, nil
				}

				iteratorFinished = true
				if inB {
					startWithItem = true
				}
				if !inA {
					return *new(T), false, err
				}

				return item, true, nil
			}
		}

		startIterator := func() bool {
			if !iter.Next() {
				return false
			}
			item, err = iter.Get()
			if err != nil {
				return false
			}
			_, _, _, fnErr := fn(count, item)
			if fnErr != nil {
				err = fnErr
				return false
			}

			startWithItem = true
			curr = FromFunc(next())
			return true
		}

		return &iterator[Iterator[T]]{
			next: func() bool {
				if finished || err != nil {
					return false
				}

				if curr == nil {
					started := startIterator()
					if started {
						count++
					}
					return started
				}

				for curr.Next() {
				}
				if err != nil {
					return false
				}
				if startWithItem {
					curr = FromFunc(next())
					return true
				}

				return startIterator()
			},
			get: func() (Iterator[T], error) {
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

// Split returns a modifier that splits the iterator when it encounters
// an item equal to sep
func Split[T comparable](sep T) Modifier[T, Iterator[T]] {
	return SplitFunc(func(_ int, item T) (bool, bool, bool, error) {
		return item == sep, false, false, nil
	})
}

// SplitLeading is like Split but includes sep at the end of leading part.
func SplitLeading[T comparable](sep T) Modifier[T, Iterator[T]] {
	return SplitFunc(func(_ int, item T) (bool, bool, bool, error) {
		return item == sep, true, false, nil
	})
}

// SplitTrailing is like Split but includes sep at the start of trailing part.
func SplitTrailing[T comparable](sep T) Modifier[T, Iterator[T]] {
	return SplitFunc(func(_ int, item T) (bool, bool, bool, error) {
		return item == sep, false, true, nil
	})
}

// Chunk returns a modifier that splits the iterator into smaller chunks
func Chunk[T any](size int) Modifier[T, Iterator[T]] {
	return func(iter Iterator[T]) Iterator[Iterator[T]] {
		return SplitFunc(func(i int, _ T) (bool, bool, bool, error) {
			if size == 0 {
				return false, false, false, nil
			}
			shouldSplit := i != 0 && i%int(size) == 0
			return shouldSplit, false, shouldSplit, nil
		})(iter)
	}
}

// RunsFunc is a modifier that splits the iterator into multiple iterators each containing
// matching consecutive items using fn to get a comparable key for each item.
func RunsFunc[T any, S comparable](fn func(int, T) (S, error)) Modifier[T, Iterator[T]] {
	return func(iter Iterator[T]) Iterator[Iterator[T]] {
		var lastKey S
		return SplitFunc(func(i int, item T) (bool, bool, bool, error) {
			key, err := fn(i, item)
			if err != nil {
				return false, false, false, err
			}

			if i > 0 && key != lastKey {
				lastKey = key
				return true, false, true, nil
			}
			if i == 0 {
				lastKey = key
			}
			return false, false, false, nil
		})(iter)
	}
}

// Runs is a modifier that splits the iterator into multiple iterators each containing
// matching consecutive items.
func Runs[T comparable](iter Iterator[T]) Iterator[Iterator[T]] {
	return RunsFunc(func(_ int, item T) (T, error) { return item, nil })(iter)
}
