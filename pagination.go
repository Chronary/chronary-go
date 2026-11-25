package chronary

import (
	"context"
	"iter"
	"net/url"
	"strconv"
)

// PageIterator provides paginated iteration over API resources.
//
// Lists that return next_cursor (Agents.List, Calendars.List,
// Agents.ListCalendars and Events.ListAll) are walked by cursor, so iteration
// is not limited by the API's cap on offset (10,000). Other lists are walked
// by offset.
type PageIterator[T any] struct {
	fetch func(ctx context.Context, offset, limit int, cursor string) (*ListResponse[T], error)
	limit int
	// offset is where All and Collect start: the caller's params.Offset
	// (#2720). GetPage takes its own offset argument.
	offset int
	// cursorPaged: the list returns next_cursor on every page, and a nil
	// next_cursor means the last page.
	cursorPaged bool
}

// Page represents a single page of results.
type Page[T any] struct {
	Data    []T
	Total   int
	HasMore bool
	// NextCursor, on lists that have cursors, is the cursor for the next page
	// (pass it to GetPageAfter); nil on the last page and on offset-only lists.
	NextCursor *string
}

// newPageIterator creates a PageIterator for an offset-only list. start is
// the offset All and Collect begin at (the caller's params.Offset).
func newPageIterator[T any](limit, start int, fetch func(ctx context.Context, offset, limit int) (*ListResponse[T], error)) *PageIterator[T] {
	if limit <= 0 {
		limit = 50
	}
	return &PageIterator[T]{
		fetch: func(ctx context.Context, offset, limit int, _ string) (*ListResponse[T], error) {
			return fetch(ctx, offset, limit)
		},
		limit:  limit,
		offset: max(start, 0),
	}
}

// newCursorPageIterator creates a PageIterator for a list that returns
// next_cursor. fetch gets a non-empty cursor when it should page by cursor;
// setPageQuery applies it. start is the offset All and Collect begin at (the
// caller's params.Offset); later pages follow the cursor.
func newCursorPageIterator[T any](limit, start int, fetch func(ctx context.Context, offset, limit int, cursor string) (*ListResponse[T], error)) *PageIterator[T] {
	if limit <= 0 {
		limit = 50
	}
	return &PageIterator[T]{fetch: fetch, limit: limit, offset: max(start, 0), cursorPaged: true}
}

// setPageQuery sets limit, and either cursor or offset: the API rejects a
// cursor combined with a non-zero offset, so a cursor page sends no offset.
func setPageQuery(q url.Values, offset, limit int, cursor string) {
	q.Set("limit", strconv.Itoa(limit))
	if cursor != "" {
		q.Set("cursor", cursor)
		q.Del("offset")
		return
	}
	q.Set("offset", strconv.Itoa(offset))
}

func (p *PageIterator[T]) page(resp *ListResponse[T], offset int) *Page[T] {
	hasMore := offset+len(resp.Data) < resp.Total
	if p.cursorPaged {
		hasMore = resp.NextCursor != nil && *resp.NextCursor != ""
	}
	return &Page[T]{Data: resp.Data, Total: resp.Total, HasMore: hasMore, NextCursor: resp.NextCursor}
}

// GetPage fetches the page that starts at offset. The argument is used as
// given; params.Offset only sets where All and Collect start.
func (p *PageIterator[T]) GetPage(ctx context.Context, offset int) (*Page[T], error) {
	resp, err := p.fetch(ctx, offset, p.limit, "")
	if err != nil {
		return nil, err
	}
	return p.page(resp, offset), nil
}

// GetPageAfter fetches the page that starts at cursor (a previous page's
// NextCursor). Only lists that have cursors accept one: Agents.List,
// Calendars.List, Agents.ListCalendars and Events.ListAll.
func (p *PageIterator[T]) GetPageAfter(ctx context.Context, cursor string) (*Page[T], error) {
	resp, err := p.fetch(ctx, 0, p.limit, cursor)
	if err != nil {
		return nil, err
	}
	return p.page(resp, 0), nil
}

// All returns an iter.Seq2 for use with range-over-func. It starts at the
// list params' Offset (0 when unset).
//
//	for item, err := range iter.All(ctx) {
//	    if err != nil { ... }
//	    fmt.Println(item)
//	}
func (p *PageIterator[T]) All(ctx context.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		offset := p.offset
		cursor := ""
		for {
			resp, err := p.fetch(ctx, offset, p.limit, cursor)
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
			if len(resp.Data) == 0 {
				return
			}
			offset += len(resp.Data)
			if p.cursorPaged {
				if resp.NextCursor == nil || *resp.NextCursor == "" {
					return
				}
				// The cursor just sent coming back would page forever.
				if *resp.NextCursor == cursor {
					var zero T
					yield(zero, &Error{Type: "api_error", Message: "chronary: list returned the same next_cursor twice; stopping"})
					return
				}
				cursor = *resp.NextCursor
				continue
			}
			if offset >= resp.Total {
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
