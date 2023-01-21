package buffer

import (
	"sync"
)

type Buffer[T any] struct {
	buffer []T
	lock   sync.Mutex
	closed bool
}

func (b *Buffer[T]) Push(val T) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.buffer = append(b.buffer, val)
}

func (b *Buffer[T]) Pop() (T, bool) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if len(b.buffer) == 0 {
		return *new(T), false
	}
	head := b.buffer[0]
	b.buffer = b.buffer[1:]
	return head, true
}

func (b *Buffer[T]) Close() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.buffer = nil
	b.closed = true
}

func (b *Buffer[T]) IsEmpty() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return len(b.buffer) == 0
}

func (b *Buffer[T]) IsClosed() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.closed
}

func New[T any]() *Buffer[T] {
	return &Buffer[T]{}
}

type Iterator[T any] struct {
	Buffer *Buffer[T]

	nextChan  chan<- *Iterator[T]
	closeChan chan<- *Iterator[T]

	itemChan chan T
	errChan  chan error

	curr T
	err  error
}

func (iter *Iterator[T]) SendItem(val T)    { iter.itemChan <- val }
func (iter *Iterator[T]) SendErr(err error) { iter.errChan <- err }
func (iter *Iterator[T]) End()              { close(iter.itemChan) }

func (iter *Iterator[T]) Next() bool {
	if iter.Buffer.IsClosed() || iter.err != nil {
		return false
	}

	select {
	case iter.nextChan <- iter:
		return iter.waitNext()
	default:
		item, ok := iter.Buffer.Pop()
		if ok {
			iter.curr = item
			return true
		}

		iter.nextChan <- iter
		return iter.waitNext()
	}
}

func (iter *Iterator[T]) waitNext() bool {
	select {
	case val, ok := <-iter.itemChan:
		if !ok {
			iter.itemChan = nil
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

func (iter *Iterator[T]) Close() error {
	if iter.Buffer.IsClosed() || iter.err != nil {
		return iter.err
	}
	iter.closeChan <- iter
	err := <-iter.errChan
	iter.err = err
	iter.close()
	return iter.err
}

func (iter *Iterator[T]) close() {
	if iter.itemChan != nil {
		close(iter.itemChan)
		iter.itemChan = nil
	}
	if iter.errChan != nil {
		close(iter.errChan)
		iter.errChan = nil
	}
	iter.Buffer.Close()
}

func (iter *Iterator[T]) Get() (T, error) { return iter.curr, iter.err }
func (iter *Iterator[T]) Err() error      { return iter.err }

func NewIterator[T any](buffer *Buffer[T], nextChan chan<- *Iterator[T], closeChan chan<- *Iterator[T]) *Iterator[T] {
	return &Iterator[T]{
		Buffer: buffer,

		nextChan:  nextChan,
		closeChan: closeChan,

		itemChan: make(chan T),
		errChan:  make(chan error),
	}
}
