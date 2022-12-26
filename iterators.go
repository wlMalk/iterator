package iterator

type emptyIterator[T any] struct{}

type constIterator[T any] struct {
	value T
}

type sequenceIterator[T Number] struct {
	curr    T
	step    T
	asc     bool
	started bool
}

type fibonacciIterator[T Number] struct {
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

// Fibonacci returns an iterator for fibonacci numbers
func Fibonacci[T Number]() Iterator[T] {
	return &fibonacciIterator[T]{}
}

// Ascending returns an iterator of numbers from start increasing by step
func Ascending[T Number](start T, step T) Iterator[T] {
	if step == 0 {
		return Const(start)
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
func Descending[T Number](start T, step T) Iterator[T] {
	if step == 0 {
		return Const(start)
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

func (iter *fibonacciIterator[T]) Next() bool { return true }
func (iter *fibonacciIterator[T]) Get() (T, error) {
	if iter.x1 == 0 && iter.x2 == 0 {
		iter.x1 = 1
		return 0, nil
	} else if iter.x1 == 1 && iter.x2 == 0 {
		iter.x1 = 0
		iter.x2 = 1
		return 1, nil
	}
	iter.x1, iter.x2 = iter.x2, iter.x1+iter.x2
	return iter.x2, nil
}
func (iter *fibonacciIterator[T]) Close() error { return nil }
func (iter *fibonacciIterator[T]) Err() error   { return nil }
