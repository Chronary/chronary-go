package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// BillingService exposes the read-only billing summary and agent-initiated
// upgrade requests. Agents cannot pay, cancel, or change payment methods: ask
// a human to upgrade by creating an upgrade request and handing over its
// ApprovalURL.
type BillingService struct {
	service
}

// Get returns the billing summary for the org (plan, status, upgrade options,
// and any pending upgrade request). Org keys (chr_sk_) only.
func (s *BillingService) Get(ctx context.Context, opts ...RequestOption) (*BillingSummary, error) {
	var summary BillingSummary
	err := s.client.do(ctx, http.MethodGet, "/v1/billing", nil, &summary, opts...)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// CreateUpgradeRequest asks a human to upgrade the account to Pro. Give the
// returned ApprovalURL to your human; only a human can pay. If Reused is true a
// request is already pending: poll it with GetUpgradeRequest (at most once a
// minute) instead of creating another. Honors WithIdempotencyKey.
func (s *BillingService) CreateUpgradeRequest(ctx context.Context, params *CreateUpgradeRequestParams, opts ...RequestOption) (*UpgradeRequest, error) {
	if params == nil {
		params = &CreateUpgradeRequestParams{}
	}
	var req UpgradeRequest
	err := s.client.do(ctx, http.MethodPost, "/v1/billing/upgrade-requests", params, &req, opts...)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// GetUpgradeRequest retrieves an upgrade request by ID.
func (s *BillingService) GetUpgradeRequest(ctx context.Context, id string, opts ...RequestOption) (*UpgradeRequest, error) {
	var req UpgradeRequest
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/billing/upgrade-requests/%s", id), nil, &req, opts...)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// CancelUpgradeRequest withdraws a pending upgrade request.
func (s *BillingService) CancelUpgradeRequest(ctx context.Context, id string, opts ...RequestOption) (*UpgradeRequest, error) {
	var req UpgradeRequest
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/billing/upgrade-requests/%s/cancel", id), nil, &req, opts...)
	if err != nil {
		return nil, err
	}
	return &req, nil
}
