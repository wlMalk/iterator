package iterator

// Unfold generates an iterator from fn using the given initial state
func Unfold[T any, S any](state S, fn func(int, S) (T, S, bool, error)) Iterator[T] {
	var count int
	var hasMore bool
	var item T
	var err error

	return FromFunc(func() (T, bool, error) {
		item, state, hasMore, err = fn(count, state)
		if err != nil || !hasMore {
			return *new(T), false, err
		}
		count++
		return item, true, nil
	})
}
