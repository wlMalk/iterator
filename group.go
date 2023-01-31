package iterator

import (
	"github.com/wlMalk/iterator/internal/buffer"
)

// GroupFunc
func GroupFunc[T any, S comparable](fn func(int, T) (S, error)) Modifier[T, Iterator[T]] {
	return func(iter Iterator[T]) Iterator[Iterator[T]] {
		g := newGroupsHandler(iter, fn)
		go g.handle()
		return g.iterator
	}
}

// Group
func Group[T comparable](iter Iterator[T]) Iterator[Iterator[T]] {
	return GroupFunc(func(_ int, item T) (T, error) {
		return item, nil
	})(iter)
}

func uniquesFunc[T any, S comparable](fn func(int, T) (S, error), uniquesOnly bool) Modifier[T, T] {
	return func(iter Iterator[T]) Iterator[T] {
		return FilterMap(func(_ int, group Iterator[T]) (T, bool, error) {
			var item T
			var isDuplicate bool
			if _, err := Iterate(group, func(i int, groupItem T) (bool, error) {
				if i > 0 {
					isDuplicate = true
					return false, nil
				}
				item = groupItem
				return true, nil
			}); err != nil {
				return *new(T), false, err
			}
			if (uniquesOnly && isDuplicate) || (!uniquesOnly && !isDuplicate) {
				return *new(T), false, nil
			}
			return item, true, nil
		})(GroupFunc(fn)(iter))
	}
}

// DuplicatesFunc
func DuplicatesFunc[T any, S comparable](fn func(int, T) (S, error)) Modifier[T, T] {
	return uniquesFunc(fn, false)
}

// Duplicates
func Duplicates[T comparable](iter Iterator[T]) Iterator[T] {
	return DuplicatesFunc(func(_ int, item T) (T, error) {
		return item, nil
	})(iter)
}

// UniquesFunc
func UniquesFunc[T any, S comparable](fn func(int, T) (S, error)) Modifier[T, T] {
	return uniquesFunc(fn, true)
}

// Uniques
func Uniques[T comparable](iter Iterator[T]) Iterator[T] {
	return UniquesFunc(func(_ int, item T) (T, error) {
		return item, nil
	})(iter)
}

type groupsHandler[T any, S comparable] struct {
	source Iterator[T]
	fn     func(int, T) (S, error)

	parentNextChan    chan *buffer.Iterator[Iterator[T]]
	parentCloseChan   chan *buffer.Iterator[Iterator[T]]
	childrenNextChan  chan *buffer.Iterator[T]
	childrenCloseChan chan *buffer.Iterator[T]

	parentBuffer *buffer.Buffer[Iterator[T]]
	iterator     *buffer.Iterator[Iterator[T]]

	keys            map[S]int
	childrenBuffers []*buffer.Buffer[T]

	finished    bool
	err         error
	closedCount int
	count       int
}

func newGroupsHandler[T any, S comparable](source Iterator[T], fn func(int, T) (S, error)) *groupsHandler[T, S] {
	parentNextChan := make(chan *buffer.Iterator[Iterator[T]])
	parentCloseChan := make(chan *buffer.Iterator[Iterator[T]])
	childrenNextChan := make(chan *buffer.Iterator[T])
	childrenCloseChan := make(chan *buffer.Iterator[T])
	parentBuffer := buffer.New[Iterator[T]]()

	return &groupsHandler[T, S]{
		source: source,
		fn:     fn,

		parentNextChan:    parentNextChan,
		parentCloseChan:   parentCloseChan,
		childrenNextChan:  childrenNextChan,
		childrenCloseChan: childrenCloseChan,

		parentBuffer: parentBuffer,
		iterator:     buffer.NewIterator(parentBuffer, parentNextChan, parentCloseChan),

		keys: make(map[S]int),
	}
}

func (g *groupsHandler[T, S]) handle() {
	for {
		select {
		case parent := <-g.parentNextChan:
			g.nextParent(parent)
		case parent := <-g.parentCloseChan:
			parent.SendErr(nil)
		case child := <-g.childrenNextChan:
			g.nextChildren(child)
		case child := <-g.childrenCloseChan:
			child.SendErr(nil)
		}
	}
}

func (g *groupsHandler[T, S]) nextParent(parent *buffer.Iterator[Iterator[T]]) {
	item, ok := parent.Buffer.Pop()
	if ok {
		parent.SendItem(item)
		return
	}

	if g.finished {
		parent.End()
		return
	}

	if g.err != nil {
		parent.SendErr(g.err)
		return
	}

	item, hasMore, err := g.progressParent()
	if err != nil {
		g.err = err
		parent.SendErr(err)
		return
	}

	if !hasMore {
		g.finished = true
		parent.End()
		return
	}

	parent.SendItem(item)
}

func (g *groupsHandler[T, S]) nextChildren(child *buffer.Iterator[T]) {
	item, ok := child.Buffer.Pop()
	if ok {
		child.SendItem(item)
		return
	}

	if g.finished {
		child.End()
		return
	}

	if g.err != nil {
		child.SendErr(g.err)
		return
	}

	item, hasMore, err := g.progressChild(child)
	if err != nil {
		g.err = err
		child.SendErr(err)
		return
	}

	if !hasMore {
		g.finished = true
		child.End()
		return
	}

	child.SendItem(item)
}

func (g *groupsHandler[T, S]) progress() (T, S, bool, error) {
	hasMore := g.source.Next()
	if !hasMore {
		// auto close source
		closeErr := g.source.Close()
		if err := g.source.Err(); err != nil {
			return *new(T), *new(S), false, err
		}
		return *new(T), *new(S), false, closeErr
	}

	item, err := g.source.Get()
	if err != nil {
		return *new(T), *new(S), false, err
	}

	key, err := g.fn(g.count, item)
	if err != nil {
		return *new(T), *new(S), false, err
	}

	g.count++

	return item, key, true, nil
}

func (g *groupsHandler[T, S]) progressParent() (Iterator[T], bool, error) {
	for {
		item, key, hasMore, err := g.progress()
		if !hasMore || err != nil {
			return nil, false, err
		}

		keyBufferIndex, ok := g.keys[key]
		if ok {
			keyBuffer := g.childrenBuffers[keyBufferIndex]
			if !keyBuffer.IsClosed() {
				keyBuffer.Push(item)
			}
			continue
		}

		buf := buffer.New[T]()
		buf.Push(item)
		g.keys[key] = len(g.childrenBuffers)
		g.childrenBuffers = append(g.childrenBuffers, buf)
		return buffer.NewIterator(buf, g.childrenNextChan, g.childrenCloseChan), true, nil
	}
}

func (g *groupsHandler[T, S]) progressChild(child *buffer.Iterator[T]) (T, bool, error) {
	for {
		item, key, hasMore, err := g.progress()
		if !hasMore || err != nil {
			return *new(T), false, err
		}

		keyBufferIndex, ok := g.keys[key]
		if !ok {
			buf := buffer.New[T]()
			buf.Push(item)
			g.keys[key] = len(g.childrenBuffers)
			g.childrenBuffers = append(g.childrenBuffers, buf)
			g.parentBuffer.Push(buffer.NewIterator(buf, g.childrenNextChan, g.childrenCloseChan))
			continue
		}

		keyBuffer := g.childrenBuffers[keyBufferIndex]
		if keyBuffer == child.Buffer {
			return item, true, nil
		}
		if !keyBuffer.IsClosed() {
			keyBuffer.Push(item)
		}
	}
}
