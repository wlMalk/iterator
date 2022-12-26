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

// Once returns an iterator with with the given value and size 1
func Once[T any](value T) Iterator[T] {
	return Pipe(Const(value), Limit[T](1))
}

// Fibonacci returns an iterator for fibonacci numbers
func Fibonacci[T Number]() Iterator[T] {
	return &fibonacciIterator[T]{}
}

// Ascending returns an iterator of numbers from start increasing by step
func Ascending[T Number](start T, step T) Iterator[T] {
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
func Descending[T Number](start T, step T) Iterator[T] {
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
func Range[T Number](start T, end T, step T) Iterator[T] {
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
		return Pipe(Ascending(start, step), LimitFunc(func(_ uint64, item T) (bool, error) {
			return item <= end, nil
		}))
	} else {
		return Pipe(Descending(start, step), LimitFunc(func(_ uint64, item T) (bool, error) {
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
