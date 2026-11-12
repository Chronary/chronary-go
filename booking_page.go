package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// BookingPageService handles public booking page API operations.
type BookingPageService struct {
	service
}

// Create creates a new booking page for a calendar.
func (s *BookingPageService) Create(ctx context.Context, params *CreateBookingPageParams, opts ...RequestOption) (*BookingPage, error) {
	var page BookingPage
	err := s.client.do(ctx, http.MethodPost, "/v1/booking-pages", params, &page, opts...)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Get retrieves a booking page by ID.
func (s *BookingPageService) Get(ctx context.Context, id string, opts ...RequestOption) (*BookingPage, error) {
	var page BookingPage
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/booking-pages/%s", id), nil, &page, opts...)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// List returns a paginated iterator over booking pages.
func (s *BookingPageService) List(ctx context.Context, params *ListBookingPagesParams) *PageIterator[BookingPage] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[BookingPage], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[BookingPage]
		err := s.client.do(ctx, http.MethodGet, "/v1/booking-pages?"+q.Encode(), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates a booking page by ID.
func (s *BookingPageService) Update(ctx context.Context, id string, params *UpdateBookingPageParams, opts ...RequestOption) (*BookingPage, error) {
	var page BookingPage
	err := s.client.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/booking-pages/%s", id), params, &page, opts...)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Delete deletes (deactivates) a booking page by ID.
func (s *BookingPageService) Delete(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/booking-pages/%s", id), nil, opts...)
}

// ListBookings returns one page of the bookings made through a booking page,
// newest start first, with the invitee's details. Pass the returned NextCursor
// as params.Cursor for the next page (nil on the last page). Bookings made
// before this endpoint was released are not linked to their page.
func (s *BookingPageService) ListBookings(ctx context.Context, id string, params *ListBookingsParams, opts ...RequestOption) (*BookingList, error) {
	path := fmt.Sprintf("/v1/booking-pages/%s/bookings", id)
	if q := addQueryParams(params); len(q) > 0 {
		path += "?" + q.Encode()
	}
	var list BookingList
	if err := s.client.do(ctx, http.MethodGet, path, nil, &list, opts...); err != nil {
		return nil, err
	}
	return &list, nil
}

// RotateSlug replaces the booking page's public link and returns the page with
// its new Slug and BookingURL. The old /book/<slug> link stops working at once.
func (s *BookingPageService) RotateSlug(ctx context.Context, id string, opts ...RequestOption) (*BookingPage, error) {
	var page BookingPage
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/booking-pages/%s/rotate-slug", id), nil, &page, opts...)
	if err != nil {
		return nil, err
	}
	return &page, nil
}
