package chronary

import (
	"context"
	"net/http"
)

// FeedbackService handles structured feedback submissions.
type FeedbackService struct {
	service
}

// Submit posts structured feedback (bug, feature, or friction) to Chronary.
//
// Rate-limited to 25 submissions per UTC day per organization. Available
// on all plans, including free. The 26th submission returns 429 with
// Retry-After set to the seconds until the next UTC midnight; the SDK
// exposes this via Error.RetryAfter.
func (s *FeedbackService) Submit(ctx context.Context, params *SubmitFeedbackParams, opts ...RequestOption) (*FeedbackAcceptedResponse, error) {
	var resp FeedbackAcceptedResponse
	err := s.client.do(ctx, http.MethodPost, "/v1/feedback", params, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
