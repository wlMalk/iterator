package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"

	it "github.com/wlMalk/iterator"
)

var ErrInvalidInput = errors.New("json: invalid input")

type Reader[T any] struct {
	dec      *json.Decoder
	closer   func() error
	ptrElems bool
	array    bool
	started  bool
	finished bool

	curr T
	err  error
}

func (r *Reader[T]) nextToken() *rune {
	tok, err := r.dec.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			r.finished = true
			return nil
		}
		r.err = err
		return nil
	}

	if tok == nil {
		r.finished = true
		return nil
	}

	delim, ok := tok.(json.Delim)
	if !ok {
		r.err = ErrInvalidInput
		return nil
	}

	token := rune(delim)
	if token == ']' {
		r.finished = true
	}

	return &token
}

func (r *Reader[T]) UseNumber() { r.dec.UseNumber() }

func (r *Reader[T]) Next() bool {
	if r.finished || r.err != nil {
		return false
	}

	if r.array && !r.started {
		tok := r.nextToken()
		if tok == nil {
			return false
		}
		if *tok != '[' {
			r.err = ErrInvalidInput
			return false
		}
		r.started = true
	}

	if !r.dec.More() {

		if r.array {
			tok := r.nextToken()
			if tok == nil || *tok != ']' {
				r.err = ErrInvalidInput
				return false
			}
		}

		r.finished = true
		return false
	}

	var item T
	if r.ptrElems {
		r.err = r.dec.Decode(item)
	} else {
		r.err = r.dec.Decode(&item)
	}
	if r.err != nil {
		return false
	}
	r.curr = item

	return true
}

func (r *Reader[T]) Get() (T, error) {
	return r.curr, r.err
}

func (r *Reader[T]) Close() error {
	if r.closer != nil {
		return r.closer()
	}
	return nil
}

func (r *Reader[T]) Err() error { return r.err }

func Decode[T any](r io.Reader) *Reader[T] {
	dec := json.NewDecoder(r)

	ptrElems := false
	var elem T
	if reflect.TypeOf(elem).Kind() == reflect.Pointer {
		ptrElems = true
	}

	if closer, ok := r.(io.ReadCloser); ok {
		return &Reader[T]{dec: dec, ptrElems: ptrElems, closer: closer.Close}
	}
	return &Reader[T]{dec: dec, ptrElems: ptrElems}
}

func Encode[T any](w io.Writer, iter it.Iterator[T]) error {
	enc := json.NewEncoder(w)
	_, err := it.Iterate(iter, func(_ uint, item T) (bool, error) {
		if err := enc.Encode(item); err != nil {
			return false, err
		}
		return true, nil
	})
	return err
}

func DecodeArray[T any](r io.Reader) *Reader[T] {
	reader := Decode[T](r)
	reader.array = true
	return reader
}

func EncodeArray[T any](w io.Writer, iter it.Iterator[T]) error {
	if _, err := w.Write([]byte{'['}); err != nil {
		return err
	}

	_, err := it.Iterate(iter, func(index uint, item T) (bool, error) {
		if index > 0 {
			if _, err := w.Write([]byte{','}); err != nil {
				return false, err
			}
		}

		b, err := json.Marshal(item)
		if err != nil {
			return false, err
		}
		_, err = w.Write(b)
		if err != nil {
			return false, err
		}

		return true, nil
	})

	if _, err := w.Write([]byte{']'}); err != nil {
		return err
	}

	return err
}

type Iterator[T any] struct {
	it.Iterator[T]
}

func (iter Iterator[T]) MarshalJSON() ([]byte, error) {
	if iter.Iterator == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer

	if err := EncodeArray[T](&buf, iter); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (iter *Iterator[T]) UnmarshalJSON(data []byte) error {
	if iter == nil {
		return errors.New("json.Iterator: UnmarshalJSON on nil pointer")
	}

	r := bytes.NewReader(data)
	reader := DecodeArray[T](r)
	iter.Iterator = reader

	return nil
}
