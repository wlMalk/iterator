package iterator

func cycle[T any](times int, reversed bool) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		var items []T
		var curr, currTime = -1, -1
		var iterFinished bool
		return FromFunc(func() (T, bool, error) {
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
		})
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
