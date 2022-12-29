package iterator

import (
	"context"

	"golang.org/x/exp/constraints"
)

type emptyIterator[T any] struct{}

type constIterator[T any] struct {
	value T
}

type sequenceIterator[T constraints.Float | constraints.Integer] struct {
	curr    T
	step    T
	asc     bool
	started bool
}

type fibonacciIterator[T constraints.Float | constraints.Integer] struct {
	x1 T
	x2 T
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

// Fibonacci returns an iterator for fibonacci numbers
func Fibonacci[T constraints.Float | constraints.Integer]() Iterator[T] {
	return &fibonacciIterator[T]{}
}

// Ascending returns an iterator of numbers from start increasing by step
func Ascending[T constraints.Float | constraints.Integer](start T, step T) Iterator[T] {
	if step == 0 {
		return Once(start)
	}
	if step < 0 {
		panic("Ascending: step cannot be less than zero")
	}
	return &sequenceIterator[T]{
		curr: start,
		step: step,
		asc:  true,
	}
}

// Descending returns an iterator of numbers from start decreasing by step
func Descending[T constraints.Float | constraints.Integer](start T, step T) Iterator[T] {
	if step == 0 {
		return Once(start)
	}
	if step < 0 {
		panic("Descending: step cannot be less than zero")
	}
	return &sequenceIterator[T]{
		curr: start,
		step: step,
		asc:  false,
	}
}

// Range returns an iterator of numbers from start to end in increments/decrements of step
// It can include end if it matches a step increment/decrement
func Range[T constraints.Float | constraints.Integer](start T, end T, step T) Iterator[T] {
	var asc bool
	var diff T
	if end > start {
		asc = true
		diff = end - start
	} else if start > end {
		asc = false
		diff = start - end
	}

	if step == 0 || start == end || diff < step {
		return Once(start)
	}

	if asc {
		return Pipe(Ascending(start, step), LimitFunc(func(_ uint, item T) (bool, error) {
			return item <= end, nil
		}))
	} else {
		return Pipe(Descending(start, step), LimitFunc(func(_ uint, item T) (bool, error) {
			return item >= end, nil
		}))
	}
}

func (iter *constIterator[T]) Next() bool      { return true }
func (iter *constIterator[T]) Get() (T, error) { return iter.value, nil }
func (iter *constIterator[T]) Close() error    { return nil }
func (iter *constIterator[T]) Err() error      { return nil }

func (iter *emptyIterator[T]) Next() bool      { return false }
func (iter *emptyIterator[T]) Get() (T, error) { return *new(T), nil }
func (iter *emptyIterator[T]) Close() error    { return nil }
func (iter *emptyIterator[T]) Err() error      { return nil }

func (iter *sequenceIterator[T]) Next() bool {
	if !iter.started {
		iter.started = true
	} else if iter.asc {
		iter.curr = iter.curr + iter.step
	} else {
		iter.curr = iter.curr - iter.step
	}

	return true
}
func (iter *sequenceIterator[T]) Get() (T, error) {
	return iter.curr, nil
}
func (iter *sequenceIterator[T]) Close() error { return nil }
func (iter *sequenceIterator[T]) Err() error   { return nil }

func (iter *fibonacciIterator[T]) Next() bool {
	if iter.x1 == 0 && iter.x2 == 0 {
		iter.x1 = 1
		return true
	} else if iter.x1 == 1 && iter.x2 == 0 {
		iter.x1 = 0
		iter.x2 = 1
	} else {
		iter.x1, iter.x2 = iter.x2, iter.x1+iter.x2
	}
	return true
}
func (iter *fibonacciIterator[T]) Get() (T, error) {
	return iter.x2, nil
}
func (iter *fibonacciIterator[T]) Close() error { return nil }
func (iter *fibonacciIterator[T]) Err() error   { return nil }

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
	var data []T
	_, err := Iterate(iter, func(_ uint, item T) (bool, error) {
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
		_, err := Iterate(iter, func(_ uint, item T) (bool, error) {
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

	_, err := Iterate(iter, func(_ uint, item KV[K, V]) (bool, error) {
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
			defer iter.Close()
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
