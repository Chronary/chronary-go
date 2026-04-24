package chronary

import (
	"context"
	"net/http"
)

// UsageService handles usage-related API operations.
type UsageService struct {
	service
}

// Get retrieves quota and usage statistics for the current billing period.
func (s *UsageService) Get(ctx context.Context, opts ...RequestOption) (*Usage, error) {
	var usage Usage
	err := s.client.do(ctx, http.MethodGet, "/v1/usage", nil, &usage, opts...)
	if err != nil {
		return nil, err
	}
	return &usage, nil
}
