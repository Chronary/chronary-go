package chronary

import "time"

// ListResponse is the paginated response envelope for list endpoints.
type ListResponse[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// --- Agent ---

// AgentType is the type of agent.
type AgentType string

const (
	AgentTypeAI       AgentType = "ai"
	AgentTypeHuman    AgentType = "human"
	AgentTypeResource AgentType = "resource"
)

// AgentStatus is the status of an agent.
type AgentStatus string

const (
	AgentStatusActive         AgentStatus = "active"
	AgentStatusPaused         AgentStatus = "paused"
	AgentStatusDecommissioned AgentStatus = "decommissioned"
)

// Agent represents a Chronary agent resource.
type Agent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        AgentType              `json:"type"`
	Description *string                `json:"description"`
	Status      AgentStatus            `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CreateAgentParams are the parameters for creating an agent.
type CreateAgentParams struct {
	Name        string                 `json:"name"`
	Type        AgentType              `json:"type"`
	Description *string                `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateAgentParams are the parameters for updating an agent.
type UpdateAgentParams struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Status      *AgentStatus           `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ListAgentsParams are the query parameters for listing agents.
type ListAgentsParams struct {
	Type   *AgentType   `url:"type,omitempty"`
	Status *AgentStatus `url:"status,omitempty"`
	Limit  int          `url:"limit,omitempty"`
	Offset int          `url:"offset,omitempty"`
}

// --- Calendar ---

// Calendar represents a Chronary calendar resource.
type Calendar struct {
	ID        string                 `json:"id"`
	AgentID   *string                `json:"agent_id"`
	Name      string                 `json:"name"`
	Timezone  string                 `json:"timezone"`
	ICalToken string                 `json:"ical_token,omitempty"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// CreateCalendarParams are the parameters for creating a calendar.
type CreateCalendarParams struct {
	Name     string                 `json:"name"`
	Timezone string                 `json:"timezone"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateCalendarParams are the parameters for updating a calendar.
type UpdateCalendarParams struct {
	Name     *string                `json:"name,omitempty"`
	Timezone *string                `json:"timezone,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ListCalendarsParams are the query parameters for listing calendars.
type ListCalendarsParams struct {
	AgentID *string `url:"agent_id,omitempty"`
	Include *string `url:"include,omitempty"`
	Limit   int     `url:"limit,omitempty"`
	Offset  int     `url:"offset,omitempty"`
}

// --- Event ---

// EventStatus is the status of a calendar event.
type EventStatus string

const (
	EventStatusConfirmed EventStatus = "confirmed"
	EventStatusTentative EventStatus = "tentative"
	EventStatusCancelled EventStatus = "cancelled"
)

// EventSource is the source of a calendar event.
type EventSource string

const (
	EventSourceInternal    EventSource = "internal"
	EventSourceExternalICal EventSource = "external_ical"
)

// Event represents a Chronary calendar event.
type Event struct {
	ID          string                 `json:"id"`
	CalendarID  string                 `json:"calendar_id"`
	Title       string                 `json:"title"`
	Description *string                `json:"description"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	AllDay      bool                   `json:"all_day"`
	Status      EventStatus            `json:"status"`
	Source      EventSource            `json:"source"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CreateEventParams are the parameters for creating an event.
type CreateEventParams struct {
	Title       string                 `json:"title"`
	StartTime   string                 `json:"start_time"`
	EndTime     string                 `json:"end_time"`
	Description *string                `json:"description,omitempty"`
	AllDay      *bool                  `json:"all_day,omitempty"`
	Status      *EventStatus           `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateEventParams are the parameters for updating an event.
type UpdateEventParams struct {
	Title       *string                `json:"title,omitempty"`
	Description *string                `json:"description,omitempty"`
	StartTime   *string                `json:"start_time,omitempty"`
	EndTime     *string                `json:"end_time,omitempty"`
	AllDay      *bool                  `json:"all_day,omitempty"`
	Status      *EventStatus           `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ListEventsParams are the query parameters for listing events.
type ListEventsParams struct {
	CalendarID  *string      `url:"calendar_id,omitempty"`
	AgentID     *string      `url:"agent_id,omitempty"`
	StartAfter  *string      `url:"start_after,omitempty"`
	StartBefore *string      `url:"start_before,omitempty"`
	Status      *EventStatus `url:"status,omitempty"`
	Source      *EventSource `url:"source,omitempty"`
	Limit       int          `url:"limit,omitempty"`
	Offset      int          `url:"offset,omitempty"`
}

// --- Availability ---

// SlotDuration is the duration of an availability slot.
type SlotDuration string

const (
	SlotDuration15m SlotDuration = "15m"
	SlotDuration30m SlotDuration = "30m"
	SlotDuration45m SlotDuration = "45m"
	SlotDuration1h  SlotDuration = "1h"
	SlotDuration2h  SlotDuration = "2h"
)

// AvailabilitySlot represents a single time slot in an availability response.
type AvailabilitySlot struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	Available bool   `json:"available"`
}

// AvailabilityResponse is the response from availability endpoints.
type AvailabilityResponse struct {
	Slots []AvailabilitySlot `json:"slots"`
}

// AvailabilityParams are the parameters for single-resource availability queries.
type AvailabilityParams struct {
	Start        string        `url:"start"`
	End          string        `url:"end"`
	SlotDuration *SlotDuration `url:"slot_duration,omitempty"`
	IncludeBusy  *bool         `url:"include_busy,omitempty"`
}

// CrossAgentAvailabilityParams are the parameters for cross-agent availability queries.
type CrossAgentAvailabilityParams struct {
	Agents       []string      `url:"agents"`
	Start        string        `url:"start"`
	End          string        `url:"end"`
	SlotDuration *SlotDuration `url:"slot_duration,omitempty"`
	Calendars    []string      `url:"calendars,omitempty"`
	IncludeBusy  *bool         `url:"include_busy,omitempty"`
}

// --- Webhook ---

// Webhook represents a Chronary webhook subscription.
type Webhook struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

// WebhookCreated is returned from webhook creation and includes the secret.
type WebhookCreated struct {
	Webhook
	Secret string `json:"secret"`
}

// CreateWebhookParams are the parameters for creating a webhook.
type CreateWebhookParams struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// UpdateWebhookParams are the parameters for updating a webhook.
type UpdateWebhookParams struct {
	URL    *string  `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
	Active *bool    `json:"active,omitempty"`
}

// ListWebhooksParams are the query parameters for listing webhooks.
type ListWebhooksParams struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

// --- iCal Subscription ---

// ICalSubscriptionStatus is the status of an iCal subscription.
type ICalSubscriptionStatus string

const (
	ICalSubscriptionStatusActive ICalSubscriptionStatus = "active"
	ICalSubscriptionStatusError  ICalSubscriptionStatus = "error"
	ICalSubscriptionStatusPaused ICalSubscriptionStatus = "paused"
)

// ICalSubscription represents an iCal feed subscription.
type ICalSubscription struct {
	ID           string                 `json:"id"`
	AgentID      string                 `json:"agent_id"`
	CalendarID   string                 `json:"calendar_id"`
	URL          string                 `json:"url"`
	Label        *string                `json:"label"`
	Status       ICalSubscriptionStatus `json:"status"`
	LastSyncedAt *string                `json:"last_synced_at"`
	LastError    *string                `json:"last_error"`
	CreatedAt    time.Time              `json:"created_at"`
}

// CreateICalSubscriptionParams are the parameters for creating an iCal subscription.
type CreateICalSubscriptionParams struct {
	CalendarID string  `json:"calendar_id"`
	URL        string  `json:"url"`
	Label      *string `json:"label,omitempty"`
}

// UpdateICalSubscriptionParams are the parameters for updating an iCal subscription.
type UpdateICalSubscriptionParams struct {
	Label *string `json:"label,omitempty"`
	URL   *string `json:"url,omitempty"`
}

// ListICalSubscriptionsParams are the query parameters for listing iCal subscriptions.
type ListICalSubscriptionsParams struct {
	Status *ICalSubscriptionStatus `url:"status,omitempty"`
	Limit  int                     `url:"limit,omitempty"`
	Offset int                     `url:"offset,omitempty"`
}

// --- Usage ---

// UsageCounter tracks used vs limit for a resource.
type UsageCounter struct {
	Used  int  `json:"used"`
	Limit *int `json:"limit"`
}

// Usage is the response from the usage endpoint.
type Usage struct {
	PeriodStart         string       `json:"period_start"`
	PeriodEnd           string       `json:"period_end"`
	Plan                string       `json:"plan"`
	Agents              UsageCounter `json:"agents"`
	Calendars           UsageCounter `json:"calendars"`
	Events              UsageCounter `json:"events"`
	APICalls            UsageCounter `json:"api_calls"`
	Webhooks            UsageCounter `json:"webhooks"`
	AvailabilityQueries UsageCounter `json:"availability_queries"`
	ICalSubscriptions   UsageCounter `json:"ical_subscriptions"`
}

// --- Webhook Event ---

// WebhookEvent represents a parsed webhook delivery payload.
type WebhookEvent struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
