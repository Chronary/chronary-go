package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// ICalSubscriptionService handles iCal subscription API operations.
type ICalSubscriptionService struct {
	service
}

// Create creates a new iCal subscription for an agent.
func (s *ICalSubscriptionService) Create(ctx context.Context, agentID string, params *CreateICalSubscriptionParams, opts ...RequestOption) (*ICalSubscription, error) {
	var sub ICalSubscription
	err := s.client.do(ctx, http.MethodPost, fmt.Sprintf("/v1/agents/%s/ical-subscriptions", agentID), params, &sub, opts...)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// Get retrieves an iCal subscription by ID.
func (s *ICalSubscriptionService) Get(ctx context.Context, id string, opts ...RequestOption) (*ICalSubscription, error) {
	var sub ICalSubscription
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/ical-subscriptions/%s", id), nil, &sub, opts...)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// List returns a paginated iterator over an agent's iCal subscriptions.
func (s *ICalSubscriptionService) List(ctx context.Context, agentID string, params *ListICalSubscriptionsParams) *PageIterator[ICalSubscription] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[ICalSubscription], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[ICalSubscription]
		err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/agents/%s/ical-subscriptions?%s", agentID, q.Encode()), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an iCal subscription by ID.
func (s *ICalSubscriptionService) Update(ctx context.Context, id string, params *UpdateICalSubscriptionParams, opts ...RequestOption) (*ICalSubscription, error) {
	var sub ICalSubscription
	err := s.client.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/ical-subscriptions/%s", id), params, &sub, opts...)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// Delete deletes an iCal subscription by ID.
func (s *ICalSubscriptionService) Delete(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/ical-subscriptions/%s", id), nil, opts...)
}

// Sync triggers an immediate sync of an iCal subscription.
func (s *ICalSubscriptionService) Sync(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodPost, fmt.Sprintf("/v1/ical-subscriptions/%s/sync", id), nil, opts...)
}
