package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// WebhookService handles webhook-related API operations.
type WebhookService struct {
	service
}

// Create creates a new webhook subscription. The returned WebhookCreated includes the secret.
func (s *WebhookService) Create(ctx context.Context, params *CreateWebhookParams, opts ...RequestOption) (*WebhookCreated, error) {
	var wh WebhookCreated
	err := s.client.do(ctx, http.MethodPost, "/v1/webhooks", params, &wh, opts...)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

// Get retrieves a webhook by ID.
func (s *WebhookService) Get(ctx context.Context, id string, opts ...RequestOption) (*Webhook, error) {
	var wh Webhook
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/webhooks/%s", id), nil, &wh, opts...)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

// List returns a paginated iterator over webhooks.
func (s *WebhookService) List(ctx context.Context, params *ListWebhooksParams) *PageIterator[Webhook] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[Webhook], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[Webhook]
		err := s.client.do(ctx, http.MethodGet, "/v1/webhooks?"+q.Encode(), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates a webhook by ID.
func (s *WebhookService) Update(ctx context.Context, id string, params *UpdateWebhookParams, opts ...RequestOption) (*Webhook, error) {
	var wh Webhook
	err := s.client.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/webhooks/%s", id), params, &wh, opts...)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

// Delete deletes a webhook by ID.
func (s *WebhookService) Delete(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/webhooks/%s", id), nil, opts...)
}

// ListDeliveries returns a single page of delivery attempts for a webhook subscription,
// along with rollup statistics. Use Limit/Offset on params to walk pages; the response
// is not iterator-shaped because the rollup stats are tied to a specific page query.
func (s *WebhookService) ListDeliveries(ctx context.Context, webhookID string, params *ListWebhookDeliveriesParams, opts ...RequestOption) (*WebhookDeliveryListResponse, error) {
	q := addQueryParams(params)
	var resp WebhookDeliveryListResponse
	path := fmt.Sprintf("/v1/webhooks/%s/deliveries", webhookID)
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	err := s.client.do(ctx, http.MethodGet, path, nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
