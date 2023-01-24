package iterator

import (
	"math/rand"
)

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
