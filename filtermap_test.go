package iterator

import (
	"testing"
)

func TestTakeWhile(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		pred     func(int, int) (bool, error)
		expected []int
	}{
		{Range(0, 5, 1), func(_, item int) (bool, error) { return item < 3, nil }, []int{0, 1, 2}},
		{FromSlice([]int{0, 1, 2, 3, 0, 1, 2, 3}), func(_, item int) (bool, error) { return item < 3, nil }, []int{0, 1, 2}},
	}

	for i := range cases {
		checkIteratorEqual(t, TakeWhile(cases[i].pred)(cases[i].iter), cases[i].expected)
	}
}

func TestSlice(t *testing.T) {
	cases := []struct {
		iter           Iterator[int]
		from, to, step int
		expected       []int
	}{
		{Range(0, 5, 1), 2, 5, 1, []int{2, 3, 4}},
		{Range(0, 5, 1), 0, 6, 1, []int{0, 1, 2, 3, 4, 5}},
		{Range(0, 5, 1), 0, 6, 2, []int{0, 2, 4}},
		{Range(0, 5, 1), 0, -1, 1, []int{0, 1, 2, 3, 4, 5}},
		{Range(0, 5, 1), 0, -1, 2, []int{0, 2, 4}},
		{Range(0, 5, 1), 2, -1, 2, []int{2, 4}},
	}

	for i := range cases {
		checkIteratorEqual(t, Slice[int](cases[i].from, cases[i].to, cases[i].step)(cases[i].iter), cases[i].expected)
	}
}
