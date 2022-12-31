package iterator

import (
	"testing"
)

func TestJoin(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Join[int](), []int{}},
		{Join(Empty[int](), Empty[int]()), []int{}},
		{Join(Empty[int](), Range(4, 6, 1)), []int{4, 5, 6}},
		{Join(Range(4, 6, 1), Empty[int]()), []int{4, 5, 6}},
		{Join(Range(1, 3, 1), Range(4, 6, 1)), []int{1, 2, 3, 4, 5, 6}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestJoinLeading(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{JoinLeading(Empty[int]())(Range(4, 6, 1)), []int{4, 5, 6}},
		{JoinLeading(Range(4, 6, 1))(Empty[int]()), []int{4, 5, 6}},
		{JoinLeading(Range(1, 3, 1))(Range(4, 6, 1)), []int{1, 2, 3, 4, 5, 6}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestJoinTrailing(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{JoinTrailing(Empty[int]())(Range(4, 6, 1)), []int{4, 5, 6}},
		{JoinTrailing(Range(4, 6, 1))(Empty[int]()), []int{4, 5, 6}},
		{JoinTrailing(Range(4, 6, 1))(Range(1, 3, 1)), []int{1, 2, 3, 4, 5, 6}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestAppend(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Append[int]()(Range(4, 6, 1)), []int{4, 5, 6}},
		{Append(4, 5, 6)(Empty[int]()), []int{4, 5, 6}},
		{Append(4, 5, 6)(Range(1, 3, 1)), []int{1, 2, 3, 4, 5, 6}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestPrepend(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		expected []int
	}{
		{Prepend[int]()(Range(4, 6, 1)), []int{4, 5, 6}},
		{Prepend(4, 5, 6)(Empty[int]()), []int{4, 5, 6}},
		{Prepend(1, 2, 3)(Range(4, 6, 1)), []int{1, 2, 3, 4, 5, 6}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}
