package iterator

// Nwise returns a modifier that returns an iterator of iterators.
// Those iterators contain size items and starting from the first item, sliding up to the last item.
func Nwise[T any](size int) Modifier[T, Iterator[T]] {
	return func(iter Iterator[T]) Iterator[Iterator[T]] {
		return &nwiseIterator[T]{
			iter:   iter,
			buffer: make([]T, 0, size),
		}
	}
}

type nwiseIterator[T any] struct {
	iter   Iterator[T]
	buffer []T
	curr   Iterator[T]
	err    error
}

func (i *nwiseIterator[T]) Next() bool {
	if len(i.buffer) == 0 { // we did not start yet
		count := 0
		size := cap(i.buffer)
		for count < cap(i.buffer) && i.iter.Next() {
			item, err := i.iter.Get()
			if err != nil {
				return false
			}
			count++
			i.buffer = append(i.buffer, item)
		}
		if count != size {
			i.err = i.iter.Err()
			return false
		}
		i.setCurr()
		return true
	}

	if !i.iter.Next() {
		i.err = i.iter.Err()
		return false
	}

	item, err := i.iter.Get()
	if err != nil {
		return false
	}
	i.buffer = append(i.buffer[1:], item)
	i.setCurr()
	return true
}

func (i *nwiseIterator[T]) setCurr() {
	c := make([]T, len(i.buffer))
	copy(c, i.buffer)
	i.curr = FromSlice(c)
}

func (i *nwiseIterator[T]) Get() (Iterator[T], error) {
	if i.err != nil {
		return nil, i.err
	}
	return i.curr, nil
}

func (i *nwiseIterator[T]) Close() error {
	return i.iter.Close()
}

func (i *nwiseIterator[T]) Err() error { return i.err }
