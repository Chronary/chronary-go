package chronary

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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
// newest first, along with rollup statistics. Walk pages by passing the response's
// NextCursor as params.Cursor (Limit/Offset also work); the response is not
// iterator-shaped because the rollup stats are tied to a specific page query.
func (s *WebhookService) ListDeliveries(ctx context.Context, webhookID string, params *ListWebhookDeliveriesParams, opts ...RequestOption) (*WebhookDeliveryListResponse, error) {
	q := addQueryParams(params)
	// addQueryParams drops a false bool, but IsTest=false is a filter of its own.
	if params != nil && params.IsTest != nil {
		q.Set("is_test", strconv.FormatBool(*params.IsTest))
	}
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

// GetDelivery returns one delivery of a webhook subscription with its payload
// and every attempt (status, duration, response body prefix, error code),
// oldest first.
func (s *WebhookService) GetDelivery(ctx context.Context, webhookID, deliveryID string, opts ...RequestOption) (*WebhookDeliveryDetail, error) {
	var resp WebhookDeliveryDetail
	path := fmt.Sprintf("/v1/webhooks/%s/deliveries/%s", webhookID, deliveryID)
	err := s.client.do(ctx, http.MethodGet, path, nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendTest sends a signed sample payload of one event type to the endpoint.
// The delivery is marked IsTest, carries X-Chronary-Test: 1, makes one attempt
// and does not count toward the monthly quota. Poll GetDelivery with the
// returned DeliveryID for the result. Limited to 10 per minute.
func (s *WebhookService) SendTest(ctx context.Context, webhookID string, params *SendTestWebhookParams, opts ...RequestOption) (*WebhookTestSent, error) {
	var resp WebhookTestSent
	path := fmt.Sprintf("/v1/webhooks/%s/test", webhookID)
	err := s.client.do(ctx, http.MethodPost, path, params, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReplayDelivery sends a failed or dropped delivery again as a new delivery
// with its own DeliveryID (Developer or Pro). The stored payload goes out unchanged, signed
// with the endpoint's current secret, with an X-Chronary-Replay-Of header. It
// follows the normal retry schedule and counts toward the monthly delivery
// quota. The endpoint must be active. Honors WithIdempotencyKey. Limited to 30
// per minute.
func (s *WebhookService) ReplayDelivery(ctx context.Context, webhookID, deliveryID string, opts ...RequestOption) (*WebhookReplaySent, error) {
	var resp WebhookReplaySent
	path := fmt.Sprintf("/v1/webhooks/%s/deliveries/%s/replay", webhookID, deliveryID)
	err := s.client.do(ctx, http.MethodPost, path, nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// RotateSecret rotates the endpoint's signing secret (Developer or Pro) and returns the new
// secret once. It signs X-Signature from the next attempt; the replaced secret
// keeps signing X-Signature-Previous until params.ExpireOldIn has passed ("now",
// "1h", "24h" by default, or "7d"). VerifySignature accepts either header, so
// receivers can switch secrets in that window. params may be nil.
func (s *WebhookService) RotateSecret(ctx context.Context, webhookID string, params *RotateWebhookSecretParams, opts ...RequestOption) (*RotatedWebhookSecret, error) {
	if params == nil {
		params = &RotateWebhookSecretParams{}
	}
	var resp RotatedWebhookSecret
	path := fmt.Sprintf("/v1/webhooks/%s/secret/rotate", webhookID)
	err := s.client.do(ctx, http.MethodPost, path, params, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
