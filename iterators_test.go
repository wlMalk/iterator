package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkIteratorEqual[T comparable](t *testing.T, iter Iterator[T], items []T) {
	defer func() {
		iter.Close()
	}()

	i := 0
	for iter.Next() {
		require.Less(t, i, len(items))

		item, err := iter.Get()
		require.NoError(t, err)
		assert.Equal(t, items[i], item)
		i++
	}
	assert.Equal(t, len(items), i)
	require.NoError(t, iter.Err())
	require.NoError(t, iter.Close())
}

func checkIteratorEqualUnordered[T comparable](t *testing.T, iter Iterator[T], items []T) {
	slc, err := ToSlice(iter)
	require.NoError(t, err)
	assert.ElementsMatch(t, slc, items)
}

func checkChannelEqual[T comparable](t *testing.T, c <-chan ValErr[T], items []T) {
	checkIteratorEqual(t, FromValErrChannel(c), items)
}

func checkFuncEqual[T comparable](t *testing.T, fn func() (T, bool, error), items []T) {
	checkIteratorEqual(t, FromFunc(fn), items)
}

func TestEmpty(t *testing.T) {
	a := Empty[int]()
	checkIteratorEqual(t, a, []int{})
}

func TestZero(t *testing.T) {
	ints := Limit[int](3)(Zero[int]())
	checkIteratorEqual(t, ints, []int{0, 0, 0})

	strs := Limit[string](3)(Zero[string]())
	checkIteratorEqual(t, strs, []string{"", "", ""})

	bools := Limit[bool](3)(Zero[bool]())
	checkIteratorEqual(t, bools, []bool{false, false, false})

	nils := Limit[*int](3)(Zero[*int]())
	checkIteratorEqual(t, nils, []*int{nil, nil, nil})
}

func TestConst(t *testing.T) {
	a := Const(1)
	a = Limit[int](5)(a)
	checkIteratorEqual(t, a, []int{1, 1, 1, 1, 1})
}

func TestAscending(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Limit[int](5)(Ascending(0, 1)), []int{0, 1, 2, 3, 4}},
		{Ascending(0, 0), []int{0}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestDescending(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Limit[int](5)(Descending(0, 1)), []int{0, -1, -2, -3, -4}},
		{Descending(0, 0), []int{0}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestRange(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Range(0, 5, 1), []int{0, 1, 2, 3, 4, 5}},
		{Range(0, 5, 2), []int{0, 2, 4}},
		{Range(5, 0, 1), []int{5, 4, 3, 2, 1, 0}},
		{Range(5, 0, 2), []int{5, 3, 1}},
		{Range(1, 1, 1), []int{1}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestFibonacci(t *testing.T) {
	a := Fibonacci[int]()
	a = Limit[int](10)(a)
	checkIteratorEqual(t, a, []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34})
}

func TestFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	a := FromSlice(slice)
	checkIteratorEqual(t, a, slice)
}

func TestToSlice(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	a := FromSlice(s1)
	s2, err := ToSlice(a)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s2)
}

func createChannel[T any](items []T) chan T {
	c := make(chan T)
	go func() {
		for _, item := range items {
			c <- item
		}
		close(c)
	}()
	return c
}

func createFunc[T any](items []T) func() (T, bool, error) {
	var curr int
	return func() (T, bool, error) {
		if curr >= len(items) {
			return *new(T), false, nil
		}
		curr++
		return items[curr-1], true, nil
	}
}

func TestFromChannel(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	a := FromChannel(createChannel(items))
	checkIteratorEqual(t, a, []int{1, 2, 3, 4, 5})
}

func TestToChannel(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	a := FromSlice(items)
	c, _ := ToChannel(a, 1)
	checkChannelEqual(t, c, []int{1, 2, 3, 4, 5})
}

func TestFromMap(t *testing.T) {
	m := map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
		"fourth": 4,
		"fifth":  5,
	}
	items := []KV[string, int]{
		{"first", 1},
		{"second", 2},
		{"third", 3},
		{"fourth", 4},
		{"fifth", 5},
	}
	a := FromMap(m)
	checkIteratorEqualUnordered(t, a, items)
}

func TestToMap(t *testing.T) {
	m1 := map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
		"fourth": 4,
		"fifth":  5,
	}
	items := []KV[string, int]{
		{"first", 1},
		{"second", 2},
		{"third", 3},
		{"fourth", 4},
		{"fifth", 5},
	}
	a := FromSlice(items)
	m2, err := ToMap(a)
	require.NoError(t, err)
	assert.Equal(t, m1, m2)
}

func TestFromFunc(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	a := FromFunc(createFunc(items))
	checkIteratorEqual(t, a, []int{1, 2, 3, 4, 5})
}

func TestToFunc(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	a := FromSlice(items)
	fn := ToFunc(a)
	checkFuncEqual(t, fn, []int{1, 2, 3, 4, 5})
}
