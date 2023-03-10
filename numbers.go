package iterator

import (
	"math"

	"golang.org/x/exp/constraints"
)

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
		return Pipe(Ascending(start, step), TakeWhile(func(_ int, item T) (bool, error) {
			return item <= end, nil
		}))
	} else {
		return Pipe(Descending(start, step), TakeWhile(func(_ int, item T) (bool, error) {
			return item >= end, nil
		}))
	}
}

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

// Add returns a modifier to add x to items.
func Add[T constraints.Float | constraints.Integer | ~string](x T) Modifier[T, T] {
	return Map(func(_ int, item T) (T, error) { return item + x, nil })
}

// Sub returns a modifier to subtract x from items.
func Sub[T constraints.Float | constraints.Integer](x T) Modifier[T, T] {
	return Map(func(_ int, item T) (T, error) { return item - x, nil })
}

// Mul returns a modifier to multiply items by x.
func Mul[T constraints.Float | constraints.Integer](x T) Modifier[T, T] {
	return Map(func(_ int, item T) (T, error) { return item * x, nil })
}

// Div returns a modifier to divide items by x.
func Div[T constraints.Float | constraints.Integer](x T) Modifier[T, T] {
	return Map(func(_ int, item T) (T, error) { return item / x, nil })
}

// Pow returns a modifier to raise items to the xth power.
func Pow[T constraints.Float | constraints.Integer](x T) Modifier[T, T] {
	return Map(func(_ int, item T) (T, error) { return T(math.Pow(float64(item), float64(x))), nil })
}

// Sqrt is a modifier to get the square root of items.
func Sqrt[T constraints.Float | constraints.Integer](iter Iterator[T]) Iterator[T] {
	return Map(func(_ int, item T) (T, error) { return T(math.Sqrt(float64(item))), nil })(iter)
}

// Mod returns a modifier to mod items with x.
func Mod[T constraints.Integer](x T) Modifier[T, T] {
	return Map(func(_ int, item T) (T, error) { return item % x, nil })
}

// Clamp returns a modifier to clamps items within min and max inclusively.
func Clamp[T constraints.Ordered](min T, max T) Modifier[T, T] {
	return Map(func(_ int, item T) (T, error) {
		if item < min {
			return min, nil
		} else if item > max {
			return max, nil
		}
		return item, nil
	})
}

// Easing returns an iterator using an easing function with as many values as n.
func Easing[T constraints.Float](n int, fn func(float64) float64) Iterator[T] {
	return Unfold(0, func(_ int, state float64) (T, float64, bool, error) {
		if state > float64(n-1) {
			return *new(T), 0, false, nil
		}
		return T(fn(state / float64(n-1))), state + 1, true, nil
	})
}

// Normalize returns a modifier that scales down numbers to be in the range [0, 1].
func Normalize[T constraints.Float | constraints.Integer, S constraints.Float](min, max T) Modifier[T, S] {
	return Interpolate[T, S](min, max, 0, 1)
}

// Interpolate returns a modifier to map numbers into another range.
func Interpolate[T constraints.Float | constraints.Integer, S constraints.Float](start1, end1 T, start2, end2 S) Modifier[T, S] {
	return Map(func(_ int, item T) (S, error) {
		if start1 == end1 {
			return 0, nil
		}
		t := S(item-start1) / S(end1-start1)
		if end2 < start2 {
			start2, end2 = end2, start2
			t = 1.0 - t
		}
		return start2 + (end2-start2)*t, nil
	})
}
