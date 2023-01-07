package csv

import (
	"encoding/csv"
	"errors"
	"io"

	it "github.com/wlMalk/iterator"
)

type Reader[T ~string] struct {
	reader       *csv.Reader
	closer       func() error
	expectHeader bool
	started      bool
	finished     bool

	header []string
	curr   []T
	err    error
}

func (r *Reader[T]) Delimiter(delim rune)       { r.reader.Comma = delim }
func (r *Reader[T]) Comment(comment rune)       { r.reader.Comment = comment }
func (r *Reader[T]) TrimLeadingSpace(trim bool) { r.reader.TrimLeadingSpace = trim }
func (r *Reader[T]) LazyQuotes(lazyQuotes bool) { r.reader.LazyQuotes = lazyQuotes }
func (r *Reader[T]) FieldsCount(count int)      { r.reader.FieldsPerRecord = count }
func (r *Reader[T]) FixedFieldsCount()          { r.reader.FieldsPerRecord = 0 }
func (r *Reader[T]) VaryingFieldsCount()        { r.reader.FieldsPerRecord = -1 }
func (r *Reader[T]) ExpectHeader(header bool)   { r.expectHeader = header }

func (r *Reader[T]) Header() ([]string, error) {
	if !r.expectHeader {
		return nil, errors.New("csv: not expecting header")
	}
	if r.header != nil {
		return r.header, nil
	}
	if r.err != nil {
		return nil, r.err
	}
	if r.started {
		return nil, r.err
	}

	header, err := r.reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			r.finished = true
			return nil, nil
		}
		r.err = err
		return nil, err
	}
	r.header = header

	return r.header, nil
}

func (r *Reader[T]) Next() bool {
	if r.finished || r.err != nil {
		return false
	}

	if r.expectHeader && r.header == nil {
		if _, err := r.Header(); err != nil {
			return false
		}
	}

	fields, err := r.reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			r.finished = true
			return false
		}
		r.err = err
		return false
	}

	if !r.started {
		r.started = true
	}

	curr := make([]T, len(fields))
	for i := range curr {
		curr[i] = T(fields[i])
	}
	r.curr = curr

	return true
}

func (r *Reader[T]) Get() ([]T, error) {
	return r.curr, r.err
}

func (r *Reader[T]) Close() error {
	if r.closer != nil {
		return r.closer()
	}
	return nil
}

func (r *Reader[T]) Err() error { return r.err }

type Writer[T ~string] struct {
	writer *csv.Writer
	iter   it.Iterator[[]T]
	header []string
}

func (w *Writer[T]) Delimiter(delim rune)   { w.writer.Comma = delim }
func (w *Writer[T]) Header(header []string) { w.header = header }
func (w *Writer[T]) UseCRLF(useCRLF bool)   { w.writer.UseCRLF = useCRLF }

func (w *Writer[T]) Write() error {
	if len(w.header) > 0 {
		if err := w.writer.Write(w.header); err != nil {
			return err
		}
	}

	_, err := it.Iterate(w.iter, func(_ uint, row []T) (bool, error) {
		fields := make([]string, len(row))
		for i := range fields {
			fields[i] = string(row[i])
		}
		if err := w.writer.Write(fields); err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		return err
	}

	return nil
}

func Read[T ~string](r io.Reader) *Reader[T] {
	reader := csv.NewReader(r)

	if closer, ok := r.(io.ReadCloser); ok {
		return &Reader[T]{reader: reader, closer: closer.Close}
	}
	return &Reader[T]{reader: reader}
}

func Write[T ~string](w io.Writer, iter it.Iterator[[]T]) *Writer[T] {
	return &Writer[T]{
		writer: csv.NewWriter(w),
		iter:   iter,
	}
}
