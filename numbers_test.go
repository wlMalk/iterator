package iterator

import (
	"testing"
)

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

func TestEasing(t *testing.T) {
	a := Easing[float64](10, func(x float64) float64 {
		if x < 0.5 {
			return 0
		}
		return 1
	})
	checkIteratorEqual(t, a, []float64{0, 0, 0, 0, 0, 1, 1, 1, 1, 1})
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		iter     Iterator[int]
		min, max int
		expected []float64
	}{
		{FromSlice([]int{0, 5, 10, 15, 20}), 0, 20, []float64{0, 0.25, 0.5, 0.75, 1}},
	}

	for i := range cases {
		checkIteratorEqual(t, Normalize[int, float64](cases[i].min, cases[i].max)(cases[i].iter), cases[i].expected)
	}
}

func TestInterpolate(t *testing.T) {
	cases := []struct {
		iter         Iterator[int]
		start1, end1 int
		start2, end2 float64
		expected     []float64
	}{
		{FromSlice([]int{0, 5, 10, 15, 20}), 0, 20, 0, 1, []float64{0, 0.25, 0.5, 0.75, 1}},
	}

	for i := range cases {
		checkIteratorEqual(t, Interpolate(cases[i].start1, cases[i].end1, cases[i].start2, cases[i].end2)(cases[i].iter), cases[i].expected)
	}
}
