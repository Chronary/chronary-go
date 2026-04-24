package chronary

import (
	"context"
	"iter"
)

// PageIterator provides paginated iteration over API resources.
type PageIterator[T any] struct {
	fetch func(ctx context.Context, offset, limit int) (*ListResponse[T], error)
	limit int
}

// Page represents a single page of results.
type Page[T any] struct {
	Data    []T
	Total   int
	HasMore bool
}

// newPageIterator creates a PageIterator with a default limit.
func newPageIterator[T any](limit int, fetch func(ctx context.Context, offset, limit int) (*ListResponse[T], error)) *PageIterator[T] {
	if limit <= 0 {
		limit = 50
	}
	return &PageIterator[T]{fetch: fetch, limit: limit}
}

// GetPage fetches a specific page of results.
func (p *PageIterator[T]) GetPage(ctx context.Context, offset int) (*Page[T], error) {
	resp, err := p.fetch(ctx, offset, p.limit)
	if err != nil {
		return nil, err
	}
	return &Page[T]{
		Data:    resp.Data,
		Total:   resp.Total,
		HasMore: offset+len(resp.Data) < resp.Total,
	}, nil
}

// All returns an iter.Seq2 for use with range-over-func.
//
//	for item, err := range iter.All(ctx) {
//	    if err != nil { ... }
//	    fmt.Println(item)
//	}
func (p *PageIterator[T]) All(ctx context.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		offset := 0
		for {
			resp, err := p.fetch(ctx, offset, p.limit)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, item := range resp.Data {
				if !yield(item, nil) {
					return
				}
			}
			offset += len(resp.Data)
			if len(resp.Data) == 0 || offset >= resp.Total {
				return
			}
		}
	}
}

// Collect fetches all pages and returns all items as a slice.
func (p *PageIterator[T]) Collect(ctx context.Context) ([]T, error) {
	var all []T
	for item, err := range p.All(ctx) {
		if err != nil {
			return nil, err
		}
		all = append(all, item)
	}
	return all, nil
}
