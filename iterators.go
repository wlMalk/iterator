package iterator

import (
	"context"
)

type emptyIterator[T any] struct{}

type constIterator[T any] struct {
	value T
}

// Empty returns an empty iterator of type T
func Empty[T any]() Iterator[T] {
	return &emptyIterator[T]{}
}

// Zero returns an infinite iterator with the zero value of type T
func Zero[T any]() Iterator[T] {
	return Const(*new(T))
}

// Const returns an infinite iterator with with the given value
func Const[T any](value T) Iterator[T] {
	return &constIterator[T]{value: value}
}

// Once returns an iterator with with the given value and size 1
func Once[T any](value T) Iterator[T] {
	return Pipe(Const(value), Limit[T](1))
}

func (iter *constIterator[T]) Next() bool      { return true }
func (iter *constIterator[T]) Get() (T, error) { return iter.value, nil }
func (iter *constIterator[T]) Close() error    { return nil }
func (iter *constIterator[T]) Err() error      { return nil }

func (iter *emptyIterator[T]) Next() bool      { return false }
func (iter *emptyIterator[T]) Get() (T, error) { return *new(T), nil }
func (iter *emptyIterator[T]) Close() error    { return nil }
func (iter *emptyIterator[T]) Err() error      { return nil }

type sliceIterator[T any] struct {
	source []T
	curr   int
}

// FromSlice returns an iterator wrapping a slice source
func FromSlice[T any](source []T) Iterator[T] {
	return &sliceIterator[T]{source: source, curr: -1}
}

// ToSlice consumes all items in the iterator into a slice
func ToSlice[T any](iter Iterator[T]) ([]T, error) {
	data := []T{}
	_, err := Iterate(iter, func(_ int, item T) (bool, error) {
		data = append(data, item)
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (iter *sliceIterator[T]) Next() bool {
	hasNext := iter.curr+1 < len(iter.source)
	if hasNext {
		iter.curr++
	}
	return hasNext
}

func (iter *sliceIterator[T]) Get() (T, error) { return iter.source[iter.curr], nil }
func (iter *sliceIterator[T]) Close() error    { return nil }
func (iter *sliceIterator[T]) Err() error      { return nil }

type (
	channelIterator[T any] struct {
		source <-chan T
		value  T
	}

	valErrChannelIterator[T any] struct {
		source <-chan ValErr[T]
		value  T
		err    error
	}
)

// FromChannel returns an iterator wrapping a channel source
func FromChannel[T any](source <-chan T) Iterator[T] {
	return &channelIterator[T]{source: source}
}

// FromValErrChannel returns an iterator wrapping a channel source with items of type ValErr
func FromValErrChannel[T any](source <-chan ValErr[T]) Iterator[T] {
	return &valErrChannelIterator[T]{source: source}
}

func (iter *channelIterator[T]) Next() bool {
	iter.value = *new(T)

	if iter.source == nil {
		return false
	}

	value, ok := <-iter.source
	if ok {
		iter.value = value
	} else {
		iter.source = nil
	}

	return ok
}

func (iter *valErrChannelIterator[T]) Next() bool {
	iter.value, iter.err = *new(T), nil

	if iter.source == nil {
		return false
	}

	value, ok := <-iter.source
	if ok {
		iter.value, iter.err = value.Val, value.Err
	} else {
		iter.source = nil
	}

	return ok
}

func (iter *channelIterator[T]) Get() (T, error) { return iter.value, nil }
func (iter *channelIterator[T]) Close() error    { return nil }
func (iter *channelIterator[T]) Err() error      { return nil }

func (iter *valErrChannelIterator[T]) Get() (T, error) { return iter.value, iter.err }
func (iter *valErrChannelIterator[T]) Close() error    { return nil }
func (iter *valErrChannelIterator[T]) Err() error      { return iter.err }

// ToChannel consumes all items in the iterator into a channel with size as capacity
func ToChannel[T any](iter Iterator[T], size int) (<-chan ValErr[T], func()) {
	stream := make(chan ValErr[T], size)
	cancel := make(chan struct{})

	go func() {
		_, err := Iterate(iter, func(_ int, item T) (bool, error) {
			select {
			case <-cancel:
				close(stream)
				return false, nil
			case stream <- ValErr[T]{Val: item}:
				return true, nil
			}
		})
		if err != nil {
			stream <- ValErr[T]{Err: err}
		}
		close(stream)
	}()

	return stream, func() { cancel <- struct{}{} }
}

// FromMap returns an iterator wrapping a map source
func FromMap[K comparable, V any](source map[K]V) Iterator[KV[K, V]] {
	c := make(chan KV[K, V], len(source)/4)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for k, v := range source {
			select {
			case <-ctx.Done():
				close(c)
				return
			case c <- KV[K, V]{Key: k, Val: v}:
				continue
			}
		}
		close(c)
	}()

	return OnClose(FromChannel(c), func() error {
		cancel()
		return nil
	})
}

// ToMap consumes all items in a KV iterator into a map
func ToMap[K comparable, V any](iter Iterator[KV[K, V]]) (map[K]V, error) {
	out := make(map[K]V)

	_, err := Iterate(iter, func(_ int, item KV[K, V]) (bool, error) {
		out[item.Key] = item.Val

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

// OnClose extends an iterator with a close callback
func OnClose[T any](iter Iterator[T], fn func() error) Iterator[T] {
	var err error
	var closed bool
	return &iterator[T]{
		next: iter.Next,
		get:  iter.Get,
		err: func() error {
			if err != nil {
				return err
			}
			return iter.Err()
		},
		close: func() error {
			if closed {
				return err
			}
			defer func() {
				closed = true
				iter.Close()
			}()
			if err = fn(); err != nil {
				return err
			}
			err = iter.Close()
			return err
		},
	}
}

type funcIterator[T any] struct {
	next func() (T, bool, error)

	value T
	err   error

	done bool
}

func (iter *funcIterator[T]) Next() bool {
	if iter.done {
		return false
	}

	var hasMore bool
	iter.value, hasMore, iter.err = iter.next()
	if !hasMore {
		iter.done = true
	}

	return hasMore
}

func (iter *funcIterator[T]) Get() (T, error) { return iter.value, iter.err }
func (iter *funcIterator[T]) Err() error      { return iter.err }
func (iter *funcIterator[T]) Close() error    { return nil }

// FromFunc returns an iterator wrapping a func source
func FromFunc[T any](next func() (T, bool, error)) Iterator[T] {
	return &funcIterator[T]{
		next: next,
	}
}

// ToFunc returns a function to consume all items in the iterator
func ToFunc[T any](iter Iterator[T]) func() (T, bool, error) {
	return func() (T, bool, error) {
		if iter.Next() {
			item, err := iter.Get()
			if err != nil {
				return *new(T), false, err
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
	}
}
