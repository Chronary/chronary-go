package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// AvailabilityService handles availability-related API operations.
type AvailabilityService struct {
	service
}

// ForAgent checks availability for a specific agent.
func (s *AvailabilityService) ForAgent(ctx context.Context, agentID string, params *AvailabilityParams, opts ...RequestOption) (*AvailabilityResponse, error) {
	q := addQueryParams(params)
	var resp AvailabilityResponse
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/agents/%s/availability?%s", agentID, q.Encode()), nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ForCalendar checks availability for a specific calendar.
func (s *AvailabilityService) ForCalendar(ctx context.Context, calendarID string, params *AvailabilityParams, opts ...RequestOption) (*AvailabilityResponse, error) {
	q := addQueryParams(params)
	var resp AvailabilityResponse
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/calendars/%s/availability?%s", calendarID, q.Encode()), nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Check performs a cross-agent availability query.
func (s *AvailabilityService) Check(ctx context.Context, params *CrossAgentAvailabilityParams, opts ...RequestOption) (*AvailabilityResponse, error) {
	q := addQueryParams(params)
	var resp AvailabilityResponse
	err := s.client.do(ctx, http.MethodGet, "/v1/availability?"+q.Encode(), nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
