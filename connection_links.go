package chronary

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ConnectionLinksService lets an agent request human-approved Google or
// Microsoft calendar setup without receiving provider credentials or data.
type ConnectionLinksService struct{ service }

func (s *ConnectionLinksService) Create(ctx context.Context, calendarID string, params *CreateConnectionLinkParams, opts ...RequestOption) (*ConnectionLink, error) {
	var link ConnectionLink
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/calendars/%s/connection-links", calendarID), params, &link, opts...)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// List returns setup links, newest first. An agent-scoped key sees only its
// own calendars' links.
func (s *ConnectionLinksService) List(ctx context.Context, params *ListConnectionLinksParams, opts ...RequestOption) (*ConnectionLinkList, error) {
	path := "/v1/connection-links"
	if params != nil {
		q := url.Values{}
		if params.Status != "" {
			q.Set("status", params.Status)
		}
		if params.CalendarID != "" {
			q.Set("calendar_id", params.CalendarID)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var list ConnectionLinkList
	if err := s.client.do(ctx, http.MethodGet, path, nil, &list, opts...); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ConnectionLinksService) Get(ctx context.Context, id string, opts ...RequestOption) (*ConnectionLink, error) {
	var link ConnectionLink
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/connection-links/%s", id), nil, &link, opts...)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *ConnectionLinksService) Cancel(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/connection-links/%s", id), nil, opts...)
}
