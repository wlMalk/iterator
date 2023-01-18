package iterator

import (
	"context"
	"sync"
)

// Merge takes multiple iterators and consumes them concurrently into a single iterator
func Merge[T any](iters ...Iterator[T]) Iterator[T] {
	c := make(chan ValErr[T], len(iters))
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(len(iters))

	for i := range iters {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			var stopped bool
			_, err := Iterate(iters[i], func(index uint, item T) (bool, error) {
				select {
				case <-ctx.Done():
					stopped = true
					return false, nil
				case c <- ValErr[T]{Val: item, Err: nil}:
				}
				return true, nil
			})
			if stopped {
				return
			}
			if err != nil {
				select {
				case <-ctx.Done():
				case c <- ValErr[T]{Val: *new(T), Err: err}:
					cancel()
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		cancel()
		close(c)
	}()

	return FromValErrChannel(c)
}
