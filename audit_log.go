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
	// Action filters by exact action name (e.g., "agent.create").
	Action string
	// ActorKeyPrefix filters by the first 20 chars of the actor API key.
	ActorKeyPrefix string
	// Cursor is an opaque pagination cursor from a previous response.
	Cursor string
	// Limit is the page size (1–200). Defaults to 50.
	Limit int
}

// List retrieves audit-log entries for the calling organization.
//
// Results are ordered by time descending and clamped to the per-tier
// retention window (Free: 7d, Pro: 90d). Only org-level API keys may
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
