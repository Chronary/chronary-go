package chronary

import (
	"context"
	"fmt"
	"net/http"
)

// AgentService handles agent-related API operations.
type AgentService struct {
	service
}

// Create registers your agent with Chronary. It creates a Chronary identity for an
// agent that already exists in your system, so it can own calendars, events, and webhooks.
func (s *AgentService) Create(ctx context.Context, params *CreateAgentParams, opts ...RequestOption) (*Agent, error) {
	var agent Agent
	err := s.client.do(ctx, http.MethodPost, "/v1/agents", params, &agent, opts...)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// Get retrieves an agent by ID.
func (s *AgentService) Get(ctx context.Context, id string, opts ...RequestOption) (*Agent, error) {
	var agent Agent
	err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/agents/%s", id), nil, &agent, opts...)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// List returns a paginated iterator over agents.
func (s *AgentService) List(ctx context.Context, params *ListAgentsParams) *PageIterator[Agent] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[Agent], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[Agent]
		err := s.client.do(ctx, http.MethodGet, "/v1/agents?"+q.Encode(), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// Update updates an agent by ID.
func (s *AgentService) Update(ctx context.Context, id string, params *UpdateAgentParams, opts ...RequestOption) (*Agent, error) {
	var agent Agent
	err := s.client.do(ctx, http.MethodPatch, fmt.Sprintf("/v1/agents/%s", id), params, &agent, opts...)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// Delete deletes an agent by ID.
func (s *AgentService) Delete(ctx context.Context, id string, opts ...RequestOption) error {
	return s.client.doNoContent(ctx, http.MethodDelete, fmt.Sprintf("/v1/agents/%s", id), nil, opts...)
}

// ListCalendars returns a paginated iterator over an agent's calendars.
func (s *AgentService) ListCalendars(ctx context.Context, agentID string, params *ListCalendarsParams) *PageIterator[Calendar] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[Calendar], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[Calendar]
		err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/agents/%s/calendars?%s", agentID, q.Encode()), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// ListEvents returns a paginated iterator over an agent's events.
func (s *AgentService) ListEvents(ctx context.Context, agentID string, params *ListEventsParams) *PageIterator[Event] {
	return newPageIterator(paramLimit(params), func(ctx context.Context, offset, limit int) (*ListResponse[Event], error) {
		q := addQueryParams(params)
		q.Set("limit", fmt.Sprintf("%d", limit))
		q.Set("offset", fmt.Sprintf("%d", offset))
		var resp ListResponse[Event]
		err := s.client.do(ctx, http.MethodGet, fmt.Sprintf("/v1/agents/%s/events?%s", agentID, q.Encode()), nil, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, nil
	})
}

// ListICalSubscriptions returns a paginated iterator over an agent's iCal subscriptions.
func (s *AgentService) ListICalSubscriptions(ctx context.Context, agentID string, params *ListICalSubscriptionsParams) *PageIterator[ICalSubscription] {
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

// paramLimit extracts limit from a params struct, returning 0 if nil or no Limit field.
func paramLimit(params interface{}) int {
	if params == nil {
		return 0
	}
	// Use a simple type switch for known param types
	switch p := params.(type) {
	case *ListAgentsParams:
		if p != nil {
			return p.Limit
		}
	case *ListCalendarsParams:
		if p != nil {
			return p.Limit
		}
	case *ListEventsParams:
		if p != nil {
			return p.Limit
		}
	case *ListWebhooksParams:
		if p != nil {
			return p.Limit
		}
	case *ListICalSubscriptionsParams:
		if p != nil {
			return p.Limit
		}
	case *ListProposalsParams:
		if p != nil {
			return p.Limit
		}
	}
	return 0
}
