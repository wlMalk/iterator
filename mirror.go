package iterator

// Mirror creates multiple synchronised iterators with the same underlying iterator
func Mirror[T any](iter Iterator[T], count int) []Iterator[T] {
	nextChan := make(chan int)
	closeChan := make(chan int)

	mirrors := make([]*mirroredIterator[T], count)

	for i := range mirrors {
		mirrors[i] = newMirroredIterator[T](i, nextChan, closeChan)
	}

	go func() {
		var finished bool
		var err error
		var closedCount int

		for {
			select {
			case index := <-nextChan:
				curr := mirrors[index]

				if curr.has() {
					curr.sendVal(curr.pop())
					continue
				}

				if finished {
					curr.end()
					continue
				}

				if err != nil {
					curr.sendErr(err)
					continue
				}

				hasMore := iter.Next()
				if !hasMore {
					if err = iter.Err(); err != nil {
						curr.sendErr(err)
						continue
					}

					finished = true
					curr.end()
					continue
				}

				val, err := iter.Get()
				if err != nil {
					curr.sendErr(err)
					continue
				}

				for i := range mirrors {
					if i == index || mirrors[i].done() {
						continue
					}
					mirrors[i].push(val)
				}

				curr.sendVal(val)

			case index := <-closeChan:
				curr := mirrors[index]
				closedCount++

				if curr.has() {
					curr.sendErr(nil)

					if closedCount == count {
						close(nextChan)
						close(closeChan)
						return
					}

					continue
				}

				err = iter.Close()
				finished = true
				curr.sendErr(err)

				if closedCount == count {
					close(nextChan)
					close(closeChan)
					return
				}
			}
		}
	}()

	iterators := make([]Iterator[T], count)
	for i := range mirrors {
		iterators[i] = mirrors[i]
	}

	return iterators
}

type mirroredIterator[T any] struct {
	index  int
	buffer []T

	nextChan  chan<- int
	closeChan chan<- int

	valChan chan T
	errChan chan error

	curr T
	err  error

	finished bool
}

func (iter *mirroredIterator[T]) push(val T) { iter.buffer = append(iter.buffer, val) }
func (iter *mirroredIterator[T]) pop() T {
	head := iter.buffer[0]
	iter.buffer = iter.buffer[1:]
	return head
}
func (iter *mirroredIterator[T]) has() bool         { return len(iter.buffer) > 0 }
func (iter *mirroredIterator[T]) done() bool        { return iter.finished }
func (iter *mirroredIterator[T]) sendVal(val T)     { iter.valChan <- val }
func (iter *mirroredIterator[T]) sendErr(err error) { iter.errChan <- err }
func (iter *mirroredIterator[T]) end()              { close(iter.valChan) }
func (iter *mirroredIterator[T]) close() {
	if iter.valChan != nil {
		close(iter.valChan)
		iter.valChan = nil
	}
	if iter.errChan != nil {
		close(iter.errChan)
		iter.errChan = nil
	}
	iter.buffer = nil
	iter.finished = true
}

func (iter *mirroredIterator[T]) Next() bool {
	if iter.finished || iter.err != nil {
		return false
	}

	iter.nextChan <- iter.index
	select {
	case val, ok := <-iter.valChan:
		if !ok {
			iter.finished = true
			iter.valChan = nil
			iter.close()
			return false
		}

		iter.curr = val
		return true
	case err := <-iter.errChan:
		iter.err = err
		iter.close()
		return false
	}
}

func (iter *mirroredIterator[T]) Get() (T, error) {
	return iter.curr, iter.err
}

func (iter *mirroredIterator[T]) Close() error {
	if iter.finished || iter.err != nil {
		return iter.err
	}
	iter.closeChan <- iter.index
	err := <-iter.errChan
	iter.err = err
	iter.close()
	return iter.err
}

func (iter *mirroredIterator[T]) Err() error {
	return iter.err
}

func newMirroredIterator[T any](index int, nextChan chan<- int, closeChan chan<- int) *mirroredIterator[T] {
	return &mirroredIterator[T]{
		index: index,

		nextChan:  nextChan,
		closeChan: closeChan,

		valChan: make(chan T),
		errChan: make(chan error),
	}
}
