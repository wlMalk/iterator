package iterator

import (
	"testing"
)

func TestInsert(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Insert(2, 2, 2)(Range(1, 4, 1)), []int{1, 2, 2, 2, 3, 4}},
		{Insert(0, 0)(Range(1, 4, 1)), []int{0, 1, 2, 3, 4}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestInject(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Inject(2, FromSlice([]int{2, 2}))(Range(1, 4, 1)), []int{1, 2, 2, 2, 3, 4}},
		{Inject(0, FromSlice([]int{0}))(Range(1, 4, 1)), []int{0, 1, 2, 3, 4}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestInterleave(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Interleave(1, Range(1, 4, 1), Range(5, 8, 1), Range(9, 12, 1)),
			[]int{1, 5, 9, 2, 6, 10, 3, 7, 11, 4, 8, 12}},
		{Interleave(2, Range(1, 4, 1), Range(5, 8, 1), Range(9, 12, 1)),
			[]int{1, 2, 5, 6, 9, 10, 3, 4, 7, 8, 11, 12}},
		{Interleave(0, Range(1, 4, 1), Range(5, 8, 1), Range(9, 12, 1)),
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}
