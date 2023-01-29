package iterator

// Repeat the items in the range [from, to) for as many times given.
// All other items outside this range are included in the iterator as well.
func Repeat[T any](from uint, to uint, times int) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var count uint
		if from > to {
			from, to = to, from
		}
		items := make([]T, 0, to-from)
		var curr, currTime = -1, 0
		var isRepeating bool
		return OnClose(FromFunc(func() (T, bool, error) {
			var hasMore bool
			var advanced bool
			if !isRepeating && count >= from && count < to {
				hasMore = iter.Next()
				advanced = true
				if !hasMore {
					isRepeating = true
				}
			}

			if isRepeating && len(items) > 0 {
				curr++
				if curr == len(items) {
					curr = 0
					currTime++
				}
				if currTime < times {
					return items[curr], true, nil
				}
			}
			isRepeating = false

			if (!advanced && iter.Next()) || (advanced && hasMore) {
				item, err := iter.Get()
				if err != nil {
					return *new(T), false, err
				}
				if count >= from && count < to {
					items = append(items, item)
				}
				count++
				if count == to {
					isRepeating = true
				}
				return item, true, nil
			}

			if err := iter.Err(); err != nil {
				return *new(T), false, err
			}

			if err := iter.Close(); err != nil {
				return *new(T), false, err
			}

			return *new(T), false, nil
		}), iter.Close)
	}
}

func cycle[T any](times int, reversed bool) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var items []T
		var curr, currTime = -1, 0
		var iterFinished bool
		return OnClose(FromFunc(func() (T, bool, error) {
			if !iterFinished {
				if iter.Next() {
					item, err := iter.Get()
					if err != nil {
						return *new(T), false, err
					}
					items = append(items, item)
					return item, true, nil
				}
				if err := iter.Err(); err != nil {
					return *new(T), false, err
				}
				if err := iter.Close(); err != nil {
					return *new(T), false, err
				}
				iterFinished = true
				currTime = 0
				if reversed {
					curr = len(items)
				}
			}

			if len(items) == 0 {
				return *new(T), false, nil
			}

			reversing := reversed && currTime%2 == 0
			finished := times > 0 &&
				currTime == times-1 &&
				((reversing && curr == 0) || (!reversing && curr == len(items)-1))

			if finished {
				return *new(T), false, nil
			}

			if reversing {
				if curr == 0 {
					reversing = false
					currTime++
				} else {
					curr--
				}
			} else {
				if curr == len(items)-1 {
					if reversed {
						reversing = true
					} else {
						curr = 0
					}
					currTime++
				} else {
					curr++
				}
			}

			return items[curr], true, nil
		}), iter.Close)
	}
}

// Cycle repeats the iterator for as many times given.
func Cycle[T any](times int) Modifier[T, T] {
	return cycle[T](times, false)
}

// Boomerang repeats the iterator for as many times given and alternates between reverse and original order.
func Boomerang[T any](times int) Modifier[T, T] {
	return cycle[T](times, true)
}
