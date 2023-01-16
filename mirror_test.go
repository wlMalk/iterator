package iterator

import (
	"sync"
	"testing"
)

func TestMirror(t *testing.T) {
	iter := Range(1, 10, 1)
	mirrors := Mirror(iter, 10)

	var wg sync.WaitGroup
	wg.Add(len(mirrors))
	for i := range mirrors {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			checkIteratorEqual(t, mirrors[i], []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		}(i)
	}
	wg.Wait()
}
