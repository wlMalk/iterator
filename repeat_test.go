package iterator

import (
	"testing"
)

func TestRepeat(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		from, to uint
		times    int
		expected []int
	}{
		{Range(1, 4, 1), 0, 1, 2, []int{1, 1, 1, 2, 3, 4}},
		{Range(1, 4, 1), 3, 4, 2, []int{1, 2, 3, 4, 4, 4}},
		{Range(1, 4, 1), 1, 2, 1, []int{1, 2, 2, 3, 4}},
		{Range(1, 4, 1), 1, 2, 2, []int{1, 2, 2, 2, 3, 4}},
		{Range(1, 4, 1), 1, 3, 2, []int{1, 2, 3, 2, 3, 2, 3, 4}},
		{Range(1, 4, 1), 0, 4, 1, []int{1, 2, 3, 4, 1, 2, 3, 4}},
		{Range(1, 4, 1), 1, 1, 2, []int{1, 2, 3, 4}},
		{Range(1, 4, 1), 0, 1, 0, []int{1, 2, 3, 4}},
		{Range(1, 4, 1), 1, 0, 1, []int{1, 1, 2, 3, 4}},
		{Range(1, 4, 1), 0, 10, 2, []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}},
		{Once(1), 0, 1, 1, []int{1, 1}},
		{Once(1), 0, 10, 1, []int{1, 1}},
		{Empty[int](), 0, 1, 1, []int{}},
	}

	for i := range cases {
		checkIteratorEqual(t, Repeat[int](cases[i].from, cases[i].to, cases[i].times)(cases[i].iter), cases[i].expected)
	}
}

func TestCycle(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		times    int
		expected []int
	}{
		{Range(1, 4, 1), 1, []int{1, 2, 3, 4, 1, 2, 3, 4}},
		{Range(1, 4, 1), 2, []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}},
		{Once(1), 3, []int{1, 1, 1, 1}},
		{Empty[int](), 4, []int{}},
	}

	for i := range cases {
		checkIteratorEqual(t, Cycle[int](cases[i].times)(cases[i].iter), cases[i].expected)
	}
}

func TestBoomerang(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		times    int
		expected []int
	}{
		{Range(1, 4, 1), 1, []int{1, 2, 3, 4, 4, 3, 2, 1}},
		{Range(1, 4, 1), 2, []int{1, 2, 3, 4, 4, 3, 2, 1, 1, 2, 3, 4}},
		{Once(1), 3, []int{1, 1, 1, 1}},
		{Empty[int](), 4, []int{}},
	}

	for i := range cases {
		checkIteratorEqual(t, Boomerang[int](cases[i].times)(cases[i].iter), cases[i].expected)
	}
}
