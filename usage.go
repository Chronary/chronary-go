package chronary

import (
	"context"
	"fmt"
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

// History retrieves monthly usage totals for the last months UTC calendar
// months (1-12), oldest first, ending with the current partial month. A
// months value of 0 uses the server default (6). It requires an
// organization-level key and is not counted toward API calls.
func (s *UsageService) History(ctx context.Context, months int, opts ...RequestOption) (*UsageHistory, error) {
	path := "/v1/usage/history"
	if months > 0 {
		path = fmt.Sprintf("%s?months=%d", path, months)
	}
	var history UsageHistory
	err := s.client.do(ctx, http.MethodGet, path, nil, &history, opts...)
	if err != nil {
		return nil, err
	}
	return &history, nil
}
