package chronary

import (
	"context"
	"net/http"
)

// PlansService handles the public plan-catalog endpoint.
type PlansService struct {
	service
}

// List fetches the public plan catalog (free, pro, scale, enterprise).
// No authentication required — usable on a client constructed with WithAnonymous().
// Responses are cached at the edge for 5 minutes.
func (s *PlansService) List(ctx context.Context, opts ...RequestOption) (*PlansListResponse, error) {
	var resp PlansListResponse
	err := s.client.do(ctx, http.MethodGet, "/v1/plans", nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
