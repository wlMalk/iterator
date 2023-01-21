package iterator

import (
	"github.com/wlMalk/iterator/internal/buffer"
)

// Mirror creates multiple synchronised iterators with the same underlying iterator
func Mirror[T any](iter Iterator[T], count int) []Iterator[T] {
	nextChan := make(chan *buffer.Iterator[T])
	closeChan := make(chan *buffer.Iterator[T])

	buffers := make([]*buffer.Buffer[T], count)
	mirrors := make([]Iterator[T], count)

	for i := range buffers {
		buffers[i] = buffer.New[T]()
		mirrors[i] = buffer.NewIterator(buffers[i], nextChan, closeChan)
	}

	go func() {
		var finished bool
		var err error
		var closedCount int

		for {
			select {
			case curr := <-nextChan:
				item, ok := curr.Buffer.Pop()
				if ok {
					curr.SendItem(item)
					continue
				}

				if finished {
					curr.End()
					continue
				}

				if err != nil {
					curr.SendErr(err)
					continue
				}

				hasMore := iter.Next()
				if !hasMore {
					if err = iter.Err(); err != nil {
						curr.SendErr(err)
						continue
					}

					finished = true
					curr.End()
					continue
				}

				val, err := iter.Get()
				if err != nil {
					curr.SendErr(err)
					continue
				}

				for _, b := range buffers {
					if b == curr.Buffer || b.IsClosed() {
						continue
					}
					b.Push(val)
				}

				curr.SendItem(val)

			case curr := <-closeChan:
				closedCount++

				if !curr.Buffer.IsEmpty() {
					curr.SendErr(nil)

					if closedCount == count {
						close(nextChan)
						close(closeChan)
						return
					}

					continue
				}

				err = iter.Close()
				finished = true
				curr.SendErr(err)

				if closedCount == count {
					close(nextChan)
					close(closeChan)
					return
				}
			}
		}
	}()

	return mirrors
}
