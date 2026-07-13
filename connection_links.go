package chronary

import (
	"context"
	"fmt"
	"net/http"
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
