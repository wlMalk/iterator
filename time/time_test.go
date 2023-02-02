package time

import (
	"testing"
	"time"

	"github.com/wlMalk/iterator"
	"github.com/wlMalk/iterator/internal/utils"
)

func checkIteratorEqual[T any](t *testing.T, iter iterator.Iterator[T], items []T) {
	utils.CheckIteratorEqual[T](t, iter, items)
}

func day(d int) time.Time {
	return time.Date(2023, 1, d, 0, 0, 0, 0, time.UTC)
}

func TestAscending(t *testing.T) {
	cases := []struct {
		iter     iterator.Iterator[time.Time]
		expected []time.Time
	}{
		{iterator.Limit[time.Time](5)(Ascending(day(1), time.Hour*24)), []time.Time{day(1), day(2), day(3), day(4), day(5)}},
		{Ascending(day(1), 0), []time.Time{day(1)}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestDescending(t *testing.T) {
	cases := []struct {
		iter     iterator.Iterator[time.Time]
		expected []time.Time
	}{
		{iterator.Limit[time.Time](5)(Descending(day(5), time.Hour*24)), []time.Time{day(5), day(4), day(3), day(2), day(1)}},
		{Descending(day(1), 0), []time.Time{day(1)}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestRange(t *testing.T) {
	cases := []struct {
		iter     iterator.Iterator[time.Time]
		expected []time.Time
	}{
		{Range(day(1), day(5), time.Hour*24), []time.Time{day(1), day(2), day(3), day(4), day(5)}},
		{Range(day(1), day(5), 2*time.Hour*24), []time.Time{day(1), day(3), day(5)}},
		{Range(day(5), day(1), time.Hour*24), []time.Time{day(5), day(4), day(3), day(2), day(1)}},
		{Range(day(5), day(1), 2*time.Hour*24), []time.Time{day(5), day(3), day(1)}},
		{Range(day(1), day(1), time.Hour*24), []time.Time{day(1)}},
	}

	for i := range cases {
		checkIteratorEqual(t, cases[i].iter, cases[i].expected)
	}
}

func TestDaysInMonth(t *testing.T) {
	checkIteratorEqual(t, DaysInMonth(1, 2023, time.UTC), []time.Time{
		day(1), day(2), day(3), day(4), day(5), day(6), day(7), day(8), day(9), day(10),
		day(11), day(12), day(13), day(14), day(15), day(16), day(17), day(18), day(19), day(20),
		day(21), day(22), day(23), day(24), day(25), day(26), day(27), day(28), day(29), day(30),
		day(31),
	})
}

func TestMonths(t *testing.T) {
	checkIteratorEqual(t, Months(), []time.Month{
		time.January,
		time.February,
		time.March,
		time.April,
		time.May,
		time.June,
		time.July,
		time.August,
		time.September,
		time.October,
		time.November,
		time.December,
	})
}

func TestWeekdays(t *testing.T) {
	checkIteratorEqual(t, Weekdays(), []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	})
}
