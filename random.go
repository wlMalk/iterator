package iterator

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

// RandomFunc returns an iterator for random numbers up to max.
// It uses fn as a random number generator in the interval [0.0,1.0).
func RandomFunc[T constraints.Float | constraints.Integer](max T, fn func() float64) Iterator[T] {
	return FromFunc(func() (T, bool, error) {
		x := T(fn() * float64(max))
		return x, true, nil
	})
}

// Random returns an iterator for random numbers up to max.
func Random[T constraints.Float | constraints.Integer](max T) Iterator[T] {
	return RandomFunc(max, rand.Float64)
}

// RandomBetweenFunc returns an iterator for random numbers from min and up to max.
// It uses fn as a random number generator in the interval [0.0,1.0).
func RandomBetweenFunc[T constraints.Float | constraints.Integer](min, max T, fn func() float64) Iterator[T] {
	return FromFunc(func() (T, bool, error) {
		x := T(fn()*float64(max-min) + float64(min))
		return x, true, nil
	})
}

// RandomBetween returns an iterator for random numbers from min and up to max.
func RandomBetween[T constraints.Float | constraints.Integer](min, max T) Iterator[T] {
	return RandomBetweenFunc(min, max, rand.Float64)
}

// SamplesFunc returns a modifier that returns random samples from the given iterator.
// It uses fn as a random number generator in the interval [0.0,1.0).
func SamplesFunc[T any](population, sample uint, fn func() float64) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		origPopulation, population, sample := population, population, sample
		var count uint
		var finished bool
		var err error

		return OnClose(FromFunc(
			func() (T, bool, error) {
				if finished {
					return *new(T), false, err
				}

				for sample > 0 && population > 0 && count < origPopulation && iter.Next() {
					count++

					shouldTake := fn() < float64(sample)/float64(population)
					population--

					if !shouldTake {
						continue
					}

					item, err := iter.Get()
					if err != nil {
						return *new(T), false, err
					}

					sample--

					return item, true, nil
				}

				finished = true
				err = iter.Err()

				return *new(T), false, err
			}), iter.Close)
	}
}

// Samples returns a modifier that returns random samples from the given iterator.
func Samples[T any](population, sample uint) Modifier[T, T] {
	return SamplesFunc[T](population, sample, rand.Float64)
}
