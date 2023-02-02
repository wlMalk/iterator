package iterator

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDistribute(t *testing.T) {
	dists, err := ToSlice(Limit[Iterator[int]](3)(Distribute[int](1)(Range(0, 9, 1))))
	require.NoError(t, err)
	require.Len(t, dists, 3)

	expects := [][]int{
		{0, 3, 6, 9},
		{1, 4, 7},
		{2, 5, 8},
	}
	var wg sync.WaitGroup
	wg.Add(len(dists))
	for i := range dists {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			checkIteratorEqual(t, dists[i], expects[i])
		}(i)
	}
	wg.Wait()
}
