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
	ID        string  `json:"id"`
	AgentID   *string `json:"agent_id"`
	Name      string  `json:"name"`
	Timezone  string  `json:"timezone"`
	ICalToken string  `json:"ical_token,omitempty"`
	// DefaultReminders is the calendar-level default reminder schedule applied to
	// events that don't set their own. Values are offsets in minutes before an
	// event's start_time (e.g. []int{10, 1440}); max 5, each 1–40320 (28 days).
	DefaultReminders []int                  `json:"default_reminders"`
	Metadata         map[string]interface{} `json:"metadata"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// CreateCalendarParams are the parameters for creating a calendar.
type CreateCalendarParams struct {
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
	// DefaultReminders sets the calendar-level default reminder schedule. Values
	// are offsets in minutes before an event's start_time (e.g. []int{10, 1440});
	// max 5, each 1–40320 (28 days). Omit to use the system default ([10]).
	DefaultReminders []int                  `json:"default_reminders,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateCalendarParams are the parameters for updating a calendar.
type UpdateCalendarParams struct {
	Name     *string `json:"name,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
	// DefaultReminders replaces the calendar-level default reminder schedule.
	// Values are offsets in minutes before an event's start_time (e.g.
	// []int{10, 1440}); max 5, each 1–40320 (28 days).
	DefaultReminders []int                  `json:"default_reminders,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
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
	EventSourceInternal     EventSource = "internal"
	EventSourceExternalICal EventSource = "external_ical"
)

// Event represents a Chronary calendar event.
type Event struct {
	ID          string      `json:"id"`
	CalendarID  string      `json:"calendar_id"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time"`
	AllDay      bool        `json:"all_day"`
	Status      EventStatus `json:"status"`
	Source      EventSource `json:"source"`
	// Reminders are offsets in minutes before start_time at which an
	// event.reminder webhook fires and a VALARM is emitted in the iCal feed
	// (e.g. []int{10, 1440}). May be null/absent when the event inherits the
	// calendar default; an empty slice means no reminders.
	Reminders []int                  `json:"reminders"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// CreateEventParams are the parameters for creating an event.
type CreateEventParams struct {
	Title       string       `json:"title"`
	StartTime   string       `json:"start_time"`
	EndTime     string       `json:"end_time"`
	Description *string      `json:"description,omitempty"`
	AllDay      *bool        `json:"all_day,omitempty"`
	Status      *EventStatus `json:"status,omitempty"`
	// Reminders are offsets in minutes before start_time (e.g. []int{10, 1440});
	// max 5, each 1–40320 (28 days). Each fires an event.reminder webhook and
	// shows as a VALARM in the iCal feed. Omit to inherit the calendar default
	// (then the system default [10]); send an empty slice for no reminders.
	Reminders []int                  `json:"reminders,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateEventParams are the parameters for updating an event.
type UpdateEventParams struct {
	Title       *string      `json:"title,omitempty"`
	Description *string      `json:"description,omitempty"`
	StartTime   *string      `json:"start_time,omitempty"`
	EndTime     *string      `json:"end_time,omitempty"`
	AllDay      *bool        `json:"all_day,omitempty"`
	Status      *EventStatus `json:"status,omitempty"`
	// Reminders replaces the event's reminder schedule. Offsets in minutes
	// before start_time (e.g. []int{10, 1440}); max 5, each 1–40320 (28 days).
	// An empty slice clears all reminders.
	Reminders []int                  `json:"reminders,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
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

// HoldsUsage holds the temporal-hold lifecycle counters for the current
// period. Informational — not gated by any plan limit. The funnel identity
// `created = confirmed + expired + active` holds, where `active` is
// derived (not stored). Counts cover all three end-of-hold paths: TTL
// expiry, manual release, and priority-bump.
type HoldsUsage struct {
	Created   int `json:"created"`
	Confirmed int `json:"confirmed"`
	Expired   int `json:"expired"`
}

// CrossCalendarQueriesUsage counts availability requests that touched more
// than one calendar in the current period. Informational — gated separately
// by the cross_calendar_availability capability, not by this counter.
type CrossCalendarQueriesUsage struct {
	Used int `json:"used"`
}

// Usage is the response from the usage endpoint.
type Usage struct {
	PeriodStart          string                    `json:"period_start"`
	PeriodEnd            string                    `json:"period_end"`
	Plan                 string                    `json:"plan"`
	Agents               UsageCounter              `json:"agents"`
	Calendars            UsageCounter              `json:"calendars"`
	Events               UsageCounter              `json:"events"`
	APICalls             UsageCounter              `json:"api_calls"`
	Webhooks             UsageCounter              `json:"webhooks"`
	AvailabilityQueries  UsageCounter              `json:"availability_queries"`
	ICalSubscriptions    UsageCounter              `json:"ical_subscriptions"`
	Holds                HoldsUsage                `json:"holds"`
	CrossCalendarQueries CrossCalendarQueriesUsage `json:"cross_calendar_queries"`
}

// --- Webhook Event ---

// WebhookEventType enumerates the event types Chronary can deliver to a webhook.
// Mirrors the server's WEBHOOK_EVENT_TYPES and the TypeScript / Python SDK unions.
type WebhookEventType string

const (
	WebhookEventAgentCreated       WebhookEventType = "agent.created"
	WebhookEventAgentUpdated       WebhookEventType = "agent.updated"
	WebhookEventEventCreated       WebhookEventType = "event.created"
	WebhookEventEventUpdated       WebhookEventType = "event.updated"
	WebhookEventEventDeleted       WebhookEventType = "event.deleted"
	WebhookEventEventStarted       WebhookEventType = "event.started"
	WebhookEventEventEnded         WebhookEventType = "event.ended"
	WebhookEventEventReminder      WebhookEventType = "event.reminder"
	WebhookEventEventHoldCreated   WebhookEventType = "event.hold_created"
	WebhookEventEventHoldExpired   WebhookEventType = "event.hold_expired"
	WebhookEventEventHoldReleased  WebhookEventType = "event.hold_released"
	WebhookEventEventHoldConfirmed WebhookEventType = "event.hold_confirmed"
	WebhookEventProposalCreated    WebhookEventType = "proposal.created"
	WebhookEventProposalResponded  WebhookEventType = "proposal.responded"
	WebhookEventProposalConfirmed  WebhookEventType = "proposal.confirmed"
	WebhookEventProposalExpired    WebhookEventType = "proposal.expired"
	WebhookEventProposalCancelled  WebhookEventType = "proposal.cancelled"
	WebhookEventWebhookDeactivated WebhookEventType = "webhook.deactivated"
)

// WebhookEventTypes is the full, ordered set of webhook event types — kept in
// lockstep with the server contract and the other SDK surfaces.
var WebhookEventTypes = []WebhookEventType{
	WebhookEventAgentCreated, WebhookEventAgentUpdated,
	WebhookEventEventCreated, WebhookEventEventUpdated, WebhookEventEventDeleted,
	WebhookEventEventStarted, WebhookEventEventEnded, WebhookEventEventReminder,
	WebhookEventEventHoldCreated, WebhookEventEventHoldExpired,
	WebhookEventEventHoldReleased, WebhookEventEventHoldConfirmed,
	WebhookEventProposalCreated, WebhookEventProposalResponded, WebhookEventProposalConfirmed,
	WebhookEventProposalExpired, WebhookEventProposalCancelled,
	WebhookEventWebhookDeactivated,
}

// WebhookEvent represents a parsed webhook delivery payload.
type WebhookEvent struct {
	Type WebhookEventType       `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// --- Calendar context ---

// AgentLiveStatus is the live-status field returned by GET /v1/calendars/{id}/context.
type AgentLiveStatus string

const (
	AgentLiveStatusIdle    AgentLiveStatus = "idle"
	AgentLiveStatusWorking AgentLiveStatus = "working"
	AgentLiveStatusWaiting AgentLiveStatus = "waiting"
	AgentLiveStatusError   AgentLiveStatus = "error"
)

// CalendarContext is the snapshot returned by GET /v1/calendars/{id}/context.
type CalendarContext struct {
	CalendarID   string          `json:"calendar_id"`
	Now          string          `json:"now"`
	AgentStatus  AgentLiveStatus `json:"agent_status"`
	CurrentEvent *Event          `json:"current_event"`
	NextEvent    *Event          `json:"next_event"`
	RecentEvents []Event         `json:"recent_events"`
	Upcoming     []Event         `json:"upcoming"`
}

// --- Availability rules ---

// WorkingHoursDay describes the working window for a single weekday.
type WorkingHoursDay struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// WorkingHours is a per-weekday map of working windows. Days that are absent are days off.
type WorkingHours struct {
	Mon *WorkingHoursDay `json:"mon,omitempty"`
	Tue *WorkingHoursDay `json:"tue,omitempty"`
	Wed *WorkingHoursDay `json:"wed,omitempty"`
	Thu *WorkingHoursDay `json:"thu,omitempty"`
	Fri *WorkingHoursDay `json:"fri,omitempty"`
	Sat *WorkingHoursDay `json:"sat,omitempty"`
	Sun *WorkingHoursDay `json:"sun,omitempty"`
}

// AvailabilityRules are the booking rules attached to a calendar.
type AvailabilityRules struct {
	ID                  string        `json:"id"`
	CalendarID          string        `json:"calendar_id"`
	BufferBeforeMinutes int           `json:"buffer_before_minutes"`
	BufferAfterMinutes  int           `json:"buffer_after_minutes"`
	WorkingHours        *WorkingHours `json:"working_hours"`
	Timezone            string        `json:"timezone"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

// SetAvailabilityRulesParams are the parameters for upserting availability rules.
type SetAvailabilityRulesParams struct {
	BufferBeforeMinutes *int          `json:"buffer_before_minutes,omitempty"`
	BufferAfterMinutes  *int          `json:"buffer_after_minutes,omitempty"`
	WorkingHours        *WorkingHours `json:"working_hours,omitempty"`
	Timezone            *string       `json:"timezone,omitempty"`
}

// --- Webhook deliveries ---

// WebhookDeliveryStatus is the delivery state of a webhook attempt.
type WebhookDeliveryStatus string

const (
	WebhookDeliveryPending   WebhookDeliveryStatus = "pending"
	WebhookDeliveryDelivered WebhookDeliveryStatus = "delivered"
	WebhookDeliveryFailed    WebhookDeliveryStatus = "failed"
)

// WebhookDelivery is a single delivery attempt for a webhook subscription.
type WebhookDelivery struct {
	ID             string                 `json:"id"`
	SubscriptionID string                 `json:"subscription_id"`
	EventType      string                 `json:"event_type"`
	Status         WebhookDeliveryStatus  `json:"status"`
	Attempts       int                    `json:"attempts"`
	LastAttemptAt  *string                `json:"last_attempt_at"`
	NextRetryAt    *string                `json:"next_retry_at"`
	CreatedAt      string                 `json:"created_at"`
	Payload        map[string]interface{} `json:"payload,omitempty"`
}

// WebhookDeliveryStats is the rollup returned alongside a deliveries list.
type WebhookDeliveryStats struct {
	Pending   int `json:"pending"`
	Delivered int `json:"delivered"`
	Failed    int `json:"failed"`
}

// WebhookDeliveryListResponse is the response from listing webhook deliveries.
type WebhookDeliveryListResponse struct {
	Data   []WebhookDelivery    `json:"data"`
	Total  int                  `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
	Stats  WebhookDeliveryStats `json:"stats"`
}

// ListWebhookDeliveriesParams are query parameters for listing webhook deliveries.
type ListWebhookDeliveriesParams struct {
	Status         *WebhookDeliveryStatus `url:"status,omitempty"`
	IncludePayload *bool                  `url:"include_payload,omitempty"`
	Limit          int                    `url:"limit,omitempty"`
	Offset         int                    `url:"offset,omitempty"`
}

// --- Scheduling proposals ---

// ProposalStatus is the lifecycle state of a scheduling proposal.
type ProposalStatus string

const (
	ProposalStatusPending   ProposalStatus = "pending"
	ProposalStatusConfirmed ProposalStatus = "confirmed"
	ProposalStatusExpired   ProposalStatus = "expired"
	ProposalStatusCancelled ProposalStatus = "cancelled"
)

// ProposalResponseAction is a participant's response to a proposal.
type ProposalResponseAction string

const (
	ProposalResponseAccept  ProposalResponseAction = "accept"
	ProposalResponseDecline ProposalResponseAction = "decline"
	ProposalResponseCounter ProposalResponseAction = "counter"
)

// ProposalSlot is a candidate time slot on a proposal.
type ProposalSlot struct {
	ID         string  `json:"id,omitempty"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	Weight     *int    `json:"weight,omitempty"`
	CalendarID *string `json:"calendar_id,omitempty"`
}

// ProposalResponse is a participant's recorded response.
type ProposalResponse struct {
	ID             string                 `json:"id"`
	AgentID        string                 `json:"agent_id"`
	Response       ProposalResponseAction `json:"response"`
	SelectedSlotID *string                `json:"selected_slot_id"`
	CounterSlots   []ProposalSlot         `json:"counter_slots"`
	Message        *string                `json:"message"`
	CreatedAt      string                 `json:"created_at"`
}

// ProposalSummary is the listing-shape proposal (no slots/responses).
type ProposalSummary struct {
	ID                  string                 `json:"id"`
	Title               string                 `json:"title"`
	Description         *string                `json:"description"`
	OrganizerAgentID    string                 `json:"organizer_agent_id"`
	ParticipantAgentIDs []string               `json:"participant_agent_ids"`
	CalendarID          string                 `json:"calendar_id"`
	Status              ProposalStatus         `json:"status"`
	ExpiresAt           *string                `json:"expires_at"`
	ResolvedSlot        *ProposalSlot          `json:"resolved_slot"`
	CreatedEventID      *string                `json:"created_event_id"`
	Metadata            map[string]interface{} `json:"metadata"`
	CreatedAt           string                 `json:"created_at"`
	UpdatedAt           string                 `json:"updated_at"`
}

// Proposal is the full proposal including slots and responses.
type Proposal struct {
	ProposalSummary
	Slots     []ProposalSlot     `json:"slots"`
	Responses []ProposalResponse `json:"responses"`
}

// CreateProposalParams are the parameters for creating a proposal.
type CreateProposalParams struct {
	Title               string                 `json:"title"`
	Description         *string                `json:"description,omitempty"`
	OrganizerAgentID    string                 `json:"organizer_agent_id"`
	ParticipantAgentIDs []string               `json:"participant_agent_ids"`
	CalendarID          string                 `json:"calendar_id"`
	Slots               []ProposalSlot         `json:"slots"`
	ExpiresAt           *string                `json:"expires_at,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// RespondToProposalParams are the parameters for responding to a proposal.
type RespondToProposalParams struct {
	AgentID        string                 `json:"agent_id"`
	Response       ProposalResponseAction `json:"response"`
	SelectedSlotID *string                `json:"selected_slot_id,omitempty"`
	CounterSlots   []ProposalSlot         `json:"counter_slots,omitempty"`
	Message        *string                `json:"message,omitempty"`
}

// ListProposalsParams are query parameters for listing proposals.
type ListProposalsParams struct {
	Status           *ProposalStatus `url:"status,omitempty"`
	OrganizerAgentID *string         `url:"organizer_agent_id,omitempty"`
	Limit            int             `url:"limit,omitempty"`
	Offset           int             `url:"offset,omitempty"`
}

// ResolveProposalResponse is the response from POST /scheduling/proposals/:id/resolve.
// Status is either "confirmed" (with ResolvedSlot) or "cancelled" (with Reason).
type ResolveProposalResponse struct {
	Status       ProposalStatus `json:"status"`
	ResolvedSlot *ProposalSlot  `json:"resolved_slot,omitempty"`
	Reason       string         `json:"reason,omitempty"`
}

// CancelProposalResponse is the response from POST /scheduling/proposals/:id/cancel.
type CancelProposalResponse struct {
	Status ProposalStatus `json:"status"`
}

// --- Scoped API keys ---

// ScopedAPIKey is an agent-scoped API key (no secret material — listing shape).
type ScopedAPIKey struct {
	ID        string  `json:"id"`
	KeyPrefix string  `json:"key_prefix"`
	AgentID   string  `json:"agent_id"`
	Label     *string `json:"label"`
	CreatedAt string  `json:"created_at"`
}

// CreatedScopedAPIKey is the response from POST /v1/keys; includes the full key string.
// The full key is shown once at creation time only.
type CreatedScopedAPIKey struct {
	ScopedAPIKey
	Key string `json:"key"`
}

// CreateScopedAPIKeyParams are the parameters for creating a scoped API key.
type CreateScopedAPIKeyParams struct {
	AgentID string  `json:"agent_id"`
	Label   *string `json:"label,omitempty"`
}

// --- Plans (public catalog) ---

// PlanLimits is the enforced caps for a plan tier.
type PlanLimits struct {
	Agents              *int `json:"agents"`
	Calendars           *int `json:"calendars"`
	Events              *int `json:"events"`
	APICalls            *int `json:"api_calls"`
	WebhookDeliveries   *int `json:"webhook_deliveries"`
	AvailabilityQueries *int `json:"availability_queries"`
	ICalSubscriptions   *int `json:"ical_subscriptions"`
	Proposals           *int `json:"proposals"`
}

// Plan is a single tier in the public plan catalog.
type Plan struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Tagline         string      `json:"tagline"`
	Price           *int        `json:"price"`
	Currency        *string     `json:"currency"`
	Limits          *PlanLimits `json:"limits"`
	DisplayFeatures []string    `json:"display_features"`
	Recommended     bool        `json:"recommended"`
	CustomPricing   bool        `json:"custom_pricing,omitempty"`
	ContactURL      string      `json:"contact_url,omitempty"`
}

// PlansListResponse is the response from GET /v1/plans.
type PlansListResponse struct {
	Plans []Plan `json:"plans"`
}

// --- Agent signup / verify ---

// AgentSignUpParams are the parameters for POST /v1/agent/sign-up.
type AgentSignUpParams struct {
	Email      string `json:"email"`
	AgentName  string `json:"agent_name"`
	TosVersion string `json:"tos_version"`
}

// AgentSignUpResponse is the response from POST /v1/agent/sign-up. Fields beyond
// Message are populated only on the new-org branch; the existing-org dedup branch
// returns just Message (no credentials, to block enumeration).
type AgentSignUpResponse struct {
	OrgID   string `json:"org_id,omitempty"`
	AgentID string `json:"agent_id,omitempty"`
	APIKey  string `json:"api_key,omitempty"`
	Message string `json:"message"`
}

// IsNewOrg returns true when the response contains credentials for a freshly
// created org. Use this to narrow the response before reading APIKey.
func (r *AgentSignUpResponse) IsNewOrg() bool {
	return r != nil && r.APIKey != ""
}

// AgentVerifyParams are the parameters for POST /v1/agent/verify.
type AgentVerifyParams struct {
	OTP string `json:"otp"`
}

// AgentVerifyResponse is the response from POST /v1/agent/verify.
type AgentVerifyResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

// --- Feedback ---

// FeedbackType categorizes a feedback submission.
type FeedbackType string

const (
	FeedbackTypeBug      FeedbackType = "bug"
	FeedbackTypeFeature  FeedbackType = "feature"
	FeedbackTypeFriction FeedbackType = "friction"
)

// SubmitFeedbackParams are the parameters for POST /v1/feedback.
type SubmitFeedbackParams struct {
	Type    FeedbackType           `json:"type"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// FeedbackAcceptedResponse is the response from POST /v1/feedback.
type FeedbackAcceptedResponse struct {
	Status string `json:"status"`
}

// --- Audit Log (#509) ---

// AuditLogEntry is a single audit-log record for a mutating API operation or
// auth lifecycle event.
type AuditLogEntry struct {
	ID string `json:"id"`
	// Action is the event identifier, e.g. "agent.create" or "auth.signin".
	Action string `json:"action"`
	// ActorKeyPrefix holds the first 20 chars of the actor API key (never the full key).
	ActorKeyPrefix *string `json:"actor_key_prefix"`
	// AgentID is set when the actor used a per-agent chr_ak_* key.
	AgentID *string `json:"agent_id"`
	// Resource contains entity IDs extracted from the path (e.g. "agt_1/cal_2").
	Resource   *string `json:"resource"`
	IP         *string `json:"ip"`
	Status     int     `json:"status"`
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	DurationMS int     `json:"duration_ms"`
	RequestID  *string `json:"request_id"`
	CreatedAt  string  `json:"created_at"`
}

// AuditLogPagination is the cursor envelope inside an audit-log list response.
type AuditLogPagination struct {
	// NextCursor is nil when there are no further pages.
	NextCursor *string `json:"next_cursor"`
}

// AuditLogListResponse is returned by AuditLogService.List.
type AuditLogListResponse struct {
	Data       []AuditLogEntry    `json:"data"`
	Pagination AuditLogPagination `json:"pagination"`
	// RetentionDays is the plan's audit-log retention window (days). Nil = unlimited.
	RetentionDays *int `json:"retention_days"`
	// RangeClamped is true when the requested from-date was outside the retention window.
	RangeClamped bool `json:"range_clamped"`
}
