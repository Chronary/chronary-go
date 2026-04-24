package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// CalendarService handles calendar-related API operations.
type CalendarService struct {
	service
}

// Create creates a new calendar (standalone, not scoped to an agent).
func (s *CalendarService) Create(ctx context.Context, params *CreateCalendarParams, opts ...RequestOption) (*Calendar, error) {
	var cal Calendar
	err := s.client.do(ctx, http.MethodPost, "/v1/calendars", params, &cal, opts...)
	if err != nil {
		return nil, err
	}
	return &cal, nil
}

// CreateForAgent creates a new calendar scoped to an agent.
func (s *CalendarService) CreateForAgent(ctx context.Context, agentID string, params *CreateCalendarParams, opts ...RequestOption) (*Calendar, error) {
	var cal Calendar
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/agents/%s/calendars", agentID), params, &cal, opts...)
	if err != nil {
		return nil, err
	}
	return &cal, nil
}

// Get retrieves a calendar by ID.
func (s *CalendarService) Get(ctx context.Context, id string, opts ...RequestOption) (*Calendar, error) {
	var cal Calendar
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/calendars/%s", id), nil, &cal, opts...)
	if err != nil {
		return nil, err
	}
	return &cal, nil
}

// List returns a paginated iterator over calendars.
func (s *CalendarService) List(ctx context.Context, params *ListCalendarsParams) *PageIterator[Calendar] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[Calendar], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[Calendar]
		err := s.client.do(ctx, http.MethodGet, "/v1/calendars?"+q.Encode(), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates a calendar by ID.
func (s *CalendarService) Update(ctx context.Context, id string, params *UpdateCalendarParams, opts ...RequestOption) (*Calendar, error) {
	var cal Calendar
	err := s.client.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/calendars/%s", id), params, &cal, opts...)
	if err != nil {
		return nil, err
	}
	return &cal, nil
}

// Delete deletes a calendar by ID.
func (s *CalendarService) Delete(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/calendars/%s", id), nil, opts...)
}
