package iterator

import (
	"sync"
)

// Distribute returns an iterator of the specified number of iterators.
// Each iterator contains a subset of the items and all of them can be
// consumed in parallel.
func Distribute[T any](buffer int) Modifier[T, Iterator[T]] {
	return func(iter Iterator[T]) Iterator[Iterator[T]] {
		return newDistributeIterator(iter, buffer)
	}
}

type distributeIterator[T any] struct {
	iter      Iterator[T]
	wg        sync.WaitGroup
	startChan chan struct{}
	closeChan chan int
	buffer    int
	chans     []chan ValErr[T]
	curr      Iterator[T]
	count     int
	done      bool
	err       error
	lock      sync.RWMutex
}

func newDistributeIterator[T any](iter Iterator[T], buffer int) *distributeIterator[T] {
	d := &distributeIterator[T]{
		iter:      iter,
		buffer:    buffer,
		startChan: make(chan struct{}),
		closeChan: make(chan int),
	}
	d.wg.Add(1)
	go d.run()
	return d
}

func (d *distributeIterator[T]) run() {
	// TODO
	// go func(){
	// 	d.wg.Wait()
	// }()
	defer d.iter.Close()
	<-d.startChan
	counter := 0

	for {
		select {
		case id := <-d.closeChan:
			d.lock.Lock()
			if d.chans[id] != nil {
				close(d.chans[id])
				d.chans[id] = nil
				d.wg.Done()
			}
			d.lock.Unlock()
		default:
			d.lock.Lock()
			if !d.iter.Next() {
				if err := d.iter.Err(); err != nil {
					d.chans[counter] <- ValErr[T]{*new(T), err}
				}
				if err := d.iter.Close(); err != nil {
					d.chans[counter] <- ValErr[T]{*new(T), err}
				}
				for i := range d.chans {
					if d.chans[i] != nil {
						close(d.chans[i])
						d.chans[i] = nil
					}
				}
				// close(d.closeChan)
			} else {
				item, err := d.iter.Get()
				if d.chans[counter] != nil {
					d.chans[counter] <- ValErr[T]{item, err}
					counter = (d.count + 1) % len(d.chans)
					d.count++
				}
			}
			d.lock.Unlock()
		}
	}
}

func (d *distributeIterator[T]) Next() bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.done {
		return false
	}

	if d.startChan != nil {
		d.startChan <- struct{}{}
		close(d.startChan)
		d.startChan = nil
	}

	id := len(d.chans)
	c := make(chan ValErr[T], d.buffer)
	d.chans = append(d.chans, c)
	d.curr = OnClose(FromValErrChannel(c), func() error {
		d.closeChan <- id
		return nil
	})
	d.wg.Add(1)

	return true
}

func (d *distributeIterator[T]) Get() (Iterator[T], error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.curr, d.err
}

func (d *distributeIterator[T]) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.done {
		return d.err
	}
	d.done = true
	d.wg.Done()
	return nil
}

func (d *distributeIterator[T]) Err() error {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.err
}
