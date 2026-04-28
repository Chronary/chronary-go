package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// KeysService handles agent-scoped API key operations.
type KeysService struct {
	service
}

// Create issues a new agent-scoped API key. The full key string is returned
// once on Key — callers must store it immediately; it cannot be re-fetched.
func (s *KeysService) Create(ctx context.Context, params *CreateScopedAPIKeyParams, opts ...RequestOption) (*CreatedScopedAPIKey, error) {
	var key CreatedScopedAPIKey
	err := s.client.do(ctx, http.MethodPost, "/v1/keys", params, &key, opts...)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// listScopedAPIKeysResponse mirrors the API envelope { "keys": [...] }.
type listScopedAPIKeysResponse struct {
	Keys []ScopedAPIKey `json:"keys"`
}

// List returns all agent-scoped API keys for the org. The list is unpaginated —
// the API returns every key in a single response.
func (s *KeysService) List(ctx context.Context, opts ...RequestOption) ([]ScopedAPIKey, error) {
	var resp listScopedAPIKeysResponse
	err := s.client.do(ctx, http.MethodGet, "/v1/keys", nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Keys, nil
}

// Delete revokes a scoped API key by ID. Idempotent: revoking an already-revoked key returns 404.
func (s *KeysService) Delete(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/keys/%s", id), nil, opts...)
}
