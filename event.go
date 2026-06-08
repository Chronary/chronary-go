package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// EventService handles event-related API operations.
type EventService struct {
	service
}

// Create creates a new event on a calendar.
func (s *EventService) Create(ctx context.Context, calendarID string, params *CreateEventParams, opts ...RequestOption) (*Event, error) {
	var event Event
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/calendars/%s/events", calendarID), params, &event, opts...)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Get retrieves an event by ID. Both calendar and event IDs are required because
// the API mounts single-event reads under /v1/calendars/{cal_id}/events/{id}.
func (s *EventService) Get(ctx context.Context, calendarID, eventID string, opts ...RequestOption) (*Event, error) {
	var event Event
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/calendars/%s/events/%s", calendarID, eventID), nil, &event, opts...)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// List returns a paginated iterator over events.
// Provide either CalendarID or AgentID in params to scope the query.
func (s *EventService) List(ctx context.Context, params *ListEventsParams) *PageIterator[Event] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[Event], error) {
		var basePath string
		if params != nil && params.CalendarID != nil {
			basePath = fmt.Sprintf("/v1/calendars/%s/events", *params.CalendarID)
		} else if params != nil && params.AgentID != nil {
			basePath = fmt.Sprintf("/v1/agents/%s/events", *params.AgentID)
		} else {
			return nil, &Error{Type: "validation_error", Message: "ListEventsParams requires CalendarID or AgentID"}
		}

		q := addQueryParams(params)
		// Remove path params from query string
		q.Del("calendar_id")
		q.Del("agent_id")
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[Event]
		err := s.client.do(ctx, http.MethodGet, basePath+"?"+q.Encode(), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an event by ID.
func (s *EventService) Update(ctx context.Context, calendarID, eventID string, params *UpdateEventParams, opts ...RequestOption) (*Event, error) {
	var event Event
	err := s.client.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/calendars/%s/events/%s", calendarID, eventID), params, &event, opts...)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Delete deletes an event by ID.
func (s *EventService) Delete(ctx context.Context, calendarID, eventID string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/calendars/%s/events/%s", calendarID, eventID), nil, opts...)
}

// GetByID retrieves an event by ID alone. The calendar is resolved internally
// from the event, so no calendar ID is required.
func (s *EventService) GetByID(ctx context.Context, eventID string, opts ...RequestOption) (*Event, error) {
	var event Event
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/events/%s", eventID), nil, &event, opts...)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// UpdateByID updates an event by ID alone. The calendar is resolved internally
// from the event, so no calendar ID is required.
func (s *EventService) UpdateByID(ctx context.Context, eventID string, params *UpdateEventParams, opts ...RequestOption) (*Event, error) {
	var event Event
	err := s.client.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/events/%s", eventID), params, &event, opts...)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// DeleteByID deletes an event by ID alone. The calendar is resolved internally
// from the event, so no calendar ID is required.
func (s *EventService) DeleteByID(ctx context.Context, eventID string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/events/%s", eventID), nil, opts...)
}

// Confirm promotes a held event (status="hold") to status="confirmed".
// Returns 409 if the event is not a hold or the hold has expired.
func (s *EventService) Confirm(ctx context.Context, eventID string, opts ...RequestOption) (*Event, error) {
	var event Event
	err := s.client.do(ctx, http.MethodPut, fmt.Sprintf("/v1/events/%s/confirm", eventID), nil, &event, opts...)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Release manually releases a held event before its TTL expires, freeing the slot.
// Returns 409 if the event is not a hold.
func (s *EventService) Release(ctx context.Context, eventID string, opts ...RequestOption) (*Event, error) {
	var event Event
	err := s.client.do(ctx, http.MethodPut, fmt.Sprintf("/v1/events/%s/release", eventID), nil, &event, opts...)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
