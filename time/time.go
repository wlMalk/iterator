package time

import (
	"time"

	"github.com/wlMalk/iterator"
)

type sequenceIterator struct {
	curr    time.Time
	step    time.Duration
	asc     bool
	started bool
}

// Ascending returns an iterator of time.Time from start increasing by step
func Ascending(start time.Time, step time.Duration) iterator.Iterator[time.Time] {
	if step == 0 {
		return iterator.Once(start)
	}
	if step < 0 {
		panic("Ascending: step cannot be less than zero")
	}
	return &sequenceIterator{
		curr: start,
		step: step,
		asc:  true,
	}
}

// Descending returns an iterator of time.Time from start decreasing by step
func Descending(start time.Time, step time.Duration) iterator.Iterator[time.Time] {
	if step == 0 {
		return iterator.Once(start)
	}
	if step < 0 {
		panic("Descending: step cannot be less than zero")
	}
	return &sequenceIterator{
		curr: start,
		step: step,
		asc:  false,
	}
}

// Range returns an iterator of time.Time from start to end in increments/decrements of step
// It can include end if it matches a step increment/decrement
func Range(start time.Time, end time.Time, step time.Duration) iterator.Iterator[time.Time] {
	var asc bool
	var diff time.Duration
	if end.After(start) {
		asc = true
		diff = end.Sub(start)
	} else if start.After(end) {
		asc = false
		diff = start.Sub(end)
	}

	if step == 0 || start == end || diff < step {
		return iterator.Once(start)
	}

	if asc {
		return iterator.Pipe(Ascending(start, step), iterator.LimitFunc(func(_ int, item time.Time) (bool, error) {
			return !end.Before(item), nil
		}))
	} else {
		return iterator.Pipe(Descending(start, step), iterator.LimitFunc(func(_ int, item time.Time) (bool, error) {
			return !item.Before(end), nil
		}))
	}
}

func (iter *sequenceIterator) Next() bool {
	if !iter.started {
		iter.started = true
		iter.curr.Month()
	} else if iter.asc {
		iter.curr = iter.curr.Add(iter.step)
	} else {
		iter.curr = iter.curr.Add(-iter.step)
	}

	return true
}
func (iter *sequenceIterator) Get() (time.Time, error) {
	return iter.curr, nil
}
func (iter *sequenceIterator) Close() error { return nil }
func (iter *sequenceIterator) Err() error   { return nil }

// DaysInMonth returns an iterator of time.Time for days in the given month of year
func DaysInMonth(month time.Month, year int, loc *time.Location) iterator.Iterator[time.Time] {
	if month < time.January || month > time.December {
		panic("DaysIn: invalid month")
	}
	daysCount := daysIn(month, year)
	dates := make([]time.Time, daysCount)
	for day := 1; day <= daysCount; day++ {
		dates[day-1] = time.Date(year, month, day, 0, 0, 0, 0, loc)
	}
	return iterator.FromSlice(dates)
}

// Months returns an iterator of time.Month values
func Months() iterator.Iterator[time.Month] {
	return iterator.FromSlice([]time.Month{
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

// Weekdays returns an iterator of time.Weekday values
func Weekdays() iterator.Iterator[time.Weekday] {
	return iterator.FromSlice([]time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	})
}

// all lines below this are copied from stdlib time package

var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

func daysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
