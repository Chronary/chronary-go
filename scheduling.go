package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// SchedulingService handles scheduling-proposal API operations.
type SchedulingService struct {
	service
}

// Create creates a new scheduling proposal.
func (s *SchedulingService) Create(ctx context.Context, params *CreateProposalParams, opts ...RequestOption) (*ProposalSummary, error) {
	var p ProposalSummary
	err := s.client.do(ctx, http.MethodPost, "/v1/scheduling/proposals", params, &p, opts...)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Get retrieves a proposal by ID, including its slots and responses.
func (s *SchedulingService) Get(ctx context.Context, id string, opts ...RequestOption) (*Proposal, error) {
	var p Proposal
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/scheduling/proposals/%s", id), nil, &p, opts...)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List returns a paginated iterator over proposals.
func (s *SchedulingService) List(ctx context.Context, params *ListProposalsParams) *PageIterator[ProposalSummary] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[ProposalSummary], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[ProposalSummary]
		err := s.client.do(ctx, http.MethodGet, "/v1/scheduling/proposals?"+q.Encode(), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Respond records a participant's response (accept / decline / counter) on a proposal.
func (s *SchedulingService) Respond(ctx context.Context, id string, params *RespondToProposalParams, opts ...RequestOption) (*ProposalResponse, error) {
	var resp ProposalResponse
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/scheduling/proposals/%s/respond", id), params, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Resolve forces resolution on a pending proposal — confirms it with the
// highest-weighted accepted slot, or cancels it when no slot has unanimous accept.
func (s *SchedulingService) Resolve(ctx context.Context, id string, opts ...RequestOption) (*ResolveProposalResponse, error) {
	var resp ResolveProposalResponse
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/scheduling/proposals/%s/resolve", id), nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Cancel cancels a pending proposal. The organizer can call this at any time before resolution.
func (s *SchedulingService) Cancel(ctx context.Context, id string, opts ...RequestOption) (*CancelProposalResponse, error) {
	var resp CancelProposalResponse
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/scheduling/proposals/%s/cancel", id), nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
