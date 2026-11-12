package chronary

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AuditLogService handles audit-log query operations.
type AuditLogService struct {
	service
}

// ListAuditLogParams are the optional query parameters for GET /v1/audit-log.
type ListAuditLogParams struct {
	// From is the lower bound (inclusive). Clamped to the retention window if older.
	From string
	// To is the upper bound (inclusive). Defaults to now.
	To string
	// Action filters by exact action name (e.g., "agent.create"). Legacy names
	// (key.create, key.delete, key.regenerate) are matched by their current names.
	Action string
	// ActorKeyPrefix filters by the first 20 chars of the actor API key.
	ActorKeyPrefix string
	// ActorType filters by who acted: "session", "org_key", "agent_key" or "system".
	ActorType string
	// ActorID filters by actor id (an API key id, or the organization id for console sessions).
	ActorID string
	// RequestID filters by the request id of the original call.
	RequestID string
	// Resource filters by an entity id the request touched, or a bare id prefix such as "cal_".
	Resource string
	// StatusClass filters by response status class: "2xx", "3xx", "4xx" or "5xx".
	StatusClass string
	// Category filters by action category (e.g. "keys"), legacy action names included.
	Category string
	// Cursor is an opaque pagination cursor from a previous response.
	Cursor string
	// Limit is the page size (1–200). Defaults to 50.
	Limit int
}

// List retrieves audit-log entries for the calling organization.
//
// Results are ordered by time descending and clamped to the per-tier
// retention window (Free: 3 days, Pro: 90 days, Custom: per contract). Only org-level API keys may
// call this endpoint — agent-scoped keys receive 403.
func (s *AuditLogService) List(ctx context.Context, params *ListAuditLogParams, opts ...RequestOption) (*AuditLogListResponse, error) {
	path := "/v1/audit-log"
	if params != nil {
		q := url.Values{}
		if params.From != "" {
			q.Set("from", params.From)
		}
		if params.To != "" {
			q.Set("to", params.To)
		}
		if params.Action != "" {
			q.Set("action", params.Action)
		}
		if params.ActorKeyPrefix != "" {
			q.Set("actor_key_prefix", params.ActorKeyPrefix)
		}
		for key, value := range map[string]string{
			"actor_type":   params.ActorType,
			"actor_id":     params.ActorID,
			"request_id":   params.RequestID,
			"resource":     params.Resource,
			"status_class": params.StatusClass,
			"category":     params.Category,
		} {
			if value != "" {
				q.Set(key, value)
			}
		}
		if params.Cursor != "" {
			q.Set("cursor", params.Cursor)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if len(q) > 0 {
			path = fmt.Sprintf("%s?%s", path, q.Encode())
		}
	}
	var result AuditLogListResponse
	err := s.client.do(ctx, http.MethodGet, path, nil, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
