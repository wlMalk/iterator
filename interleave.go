package iterator

type interleaveItem[T any] struct {
	item     T
	count    uint
	pending  bool
	finished bool
}

// InterleaveFunc returns an iterator that alternates between the given iterators.
// It switches to the next iterator when fn returns true.
func InterleaveFunc[T any](fn func(iterIndex int, currRun int, index uint, item T) (bool, error), iters ...Iterator[T]) Iterator[T] {
	if len(iters) == 0 {
		return Empty[T]()
	}

	remaining := len(iters)
	var currIndex int
	var currRun int
	var err error
	its := make([]interleaveItem[T], len(iters))

	return &iterator[T]{
		next: func() bool {
			if remaining == 0 || err != nil {
				return false
			}

			for tries := 0; tries <= len(iters); tries++ {
				if its[currIndex].pending {
					currRun++
					its[currIndex].count++
					its[currIndex].pending = false
					return true
				}
				if !its[currIndex].finished && iters[currIndex].Next() {
					var item T
					item, err = iters[currIndex].Get()
					if err != nil {
						return false
					}
					its[currIndex].item = item
					its[currIndex].pending = true

					if remaining == 1 {
						its[currIndex].pending = false
						return true
					}

					var shouldInterleave bool
					shouldInterleave, err = fn(currIndex, currRun, its[currIndex].count, item)
					if err != nil {
						return false
					}

					if !shouldInterleave {
						currRun++
						its[currIndex].count++
						its[currIndex].pending = false
						return true
					}
				} else if !its[currIndex].finished {
					remaining--
					its[currIndex].finished = true
				}

				currIndex = (currIndex + 1) % len(iters)
				currRun = 0
			}

			return false
		},
		get: func() (T, error) {
			if err != nil {
				return *new(T), err
			}

			return its[currIndex].item, nil
		},
		close: func() error {
			var err error
			for i := range iters {
				if closeErr := iters[i].Close(); err == nil && closeErr != nil {
					err = closeErr
				}
			}
			return err
		},
		err: func() error {
			if err != nil {
				return err
			}
			for i := range iters {
				if err := iters[i].Err(); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// InjectFunc returns a modifier that injects items from the given iterator once fn returns true
func InjectFunc[T any](in Iterator[T], fn func(uint, T) (bool, error)) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		return InterleaveFunc(func(iterIndex int, _ int, index uint, item T) (bool, error) {
			if iterIndex == 1 {
				return false, nil
			}
			return fn(index, item)
		}, iter, in)
	}
}

// Inject returns a modifier that injects items from the given iterator at the given index
func Inject[T any](index uint, in Iterator[T]) Modifier[T, T] {
	return InjectFunc(in, func(i uint, _ T) (bool, error) {
		return i == index, nil
	})
}

// InsertFunc returns a modifier that inserts the given items once fn returns true
func InsertFunc[T any](fn func(uint, T) (bool, error), items ...T) Modifier[T, T] {
	return InjectFunc(FromSlice(items), fn)
}

// Insert returns a modifier that inserts the given items once fn returns true
func Insert[T any](index uint, items ...T) Modifier[T, T] {
	return Inject(index, FromSlice(items))
}

// Interleave returns an iterator that alternates between the given iterators.
// It will try to take as many as count from each iterator in each round
func Interleave[T any](count int, iters ...Iterator[T]) Iterator[T] {
	return InterleaveFunc(func(_ int, currRun int, _ uint, _ T) (bool, error) {
		return currRun == count, nil
	}, iters...)
}
