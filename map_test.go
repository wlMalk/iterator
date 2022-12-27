package iterator

import (
	"testing"
)

func TestReplace(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Pipe(Ascending(0, 1), Limit[int](5), Replace(1, 10, 0)), []int{0, 1, 2, 3, 4}},
		{Pipe(Ascending(0, 1), Limit[int](5), Replace(1, 10, 1)), []int{0, 10, 2, 3, 4}},
		{Pipe(Const(1), Limit[int](5), Replace(1, 10, 1)), []int{10, 1, 1, 1, 1}},
		{Pipe(Empty[int](), Replace(1, 10, 1), Limit[int](5)), []int{}},
		{Pipe(Const(5), Replace(1, 10, 0), Limit[int](5)), []int{5, 5, 5, 5, 5}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestReplaceAll(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Pipe(Ascending(0, 1), Limit[int](5), ReplaceAll(1, 10)), []int{0, 10, 2, 3, 4}},
		{Pipe(Const(1), Limit[int](5), ReplaceAll(1, 10)), []int{10, 10, 10, 10, 10}},
		{Pipe(Const(5), ReplaceAll(1, 10), Limit[int](5)), []int{5, 5, 5, 5, 5}},
		{Pipe(Empty[int](), ReplaceAll(1, 10), Limit[int](5)), []int{}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}
