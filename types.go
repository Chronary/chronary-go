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
	OrgID       string                 `json:"orgId"`
	Name        string                 `json:"name"`
	Type        AgentType              `json:"type"`
	Description *string                `json:"description"`
	Status      AgentStatus            `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	// CalendarCount and ScopedKeyCount are set only by List with IncludeCounts.
	CalendarCount  *int `json:"calendar_count,omitempty"`
	ScopedKeyCount *int `json:"scoped_key_count,omitempty"`
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
// Decommissioned (deleted) agents are left out unless IncludeDecommissioned
// is true or Status is set.
type ListAgentsParams struct {
	Type   *AgentType   `url:"type,omitempty"`
	Status *AgentStatus `url:"status,omitempty"`
	// Q matches a case-insensitive substring of the agent name, or an exact agent ID.
	Q string `url:"q,omitempty"`
	// Sort is "created_at" (default, newest first) or "name" (A to Z).
	Sort                  string `url:"sort,omitempty"`
	IncludeDecommissioned bool   `url:"include_decommissioned,omitempty"`
	// IncludeCounts adds CalendarCount and ScopedKeyCount to each agent.
	IncludeCounts bool `url:"include_counts,omitempty"`
	Limit         int  `url:"limit,omitempty"`
	Offset        int  `url:"offset,omitempty"`
}

// --- Calendar ---

// Calendar represents a Chronary calendar resource.
type Calendar struct {
	ID       string  `json:"id"`
	OrgID    string  `json:"orgId"`
	AgentID  *string `json:"agentId"`
	Name     string  `json:"name"`
	Timezone string  `json:"timezone"`
	// AgentStatus is the live agent status surfaced on the calendar
	// (one of "idle", "working", "waiting", "error").
	AgentStatus string  `json:"agent_status"`
	ICalURL     string  `json:"ical_url"`
	ExternalID  *string `json:"externalId"`
	Provider    *string `json:"provider"`
	// DefaultReminders is the calendar-level default reminder schedule applied to
	// events that don't set their own. Values are offsets in minutes before an
	// event's start_time (e.g. []int{10, 1440}); max 5, each 1–40320 (28 days).
	DefaultReminders []int                  `json:"default_reminders"`
	Metadata         map[string]interface{} `json:"metadata"`
	DeletedAt        *string                `json:"deletedAt"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt"`
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
	// Q is a case-insensitive substring of the calendar name, or an exact calendar ID.
	Q      *string `url:"q,omitempty"`
	Limit  int     `url:"limit,omitempty"`
	Offset int     `url:"offset,omitempty"`
}

// --- Booking page ---

// BookingPage is a public, hosted scheduling link for a calendar.
type BookingPage struct {
	ID          string  `json:"id"`
	OrgID       string  `json:"org_id"`
	CalendarID  string  `json:"calendar_id"`
	AgentID     *string `json:"agent_id"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	// DurationMinutes is the slot length in minutes.
	DurationMinutes  int    `json:"duration_minutes"`
	BufferMinutes    int    `json:"buffer_minutes"`
	WindowDays       int    `json:"window_days"`
	MinNoticeMinutes int    `json:"min_notice_minutes"`
	Timezone         string `json:"timezone"`
	// AvailabilityConstraints is the weekly working-hours restriction; nil = any time in the window.
	AvailabilityConstraints *WorkingHours `json:"availability_constraints"`
	Active                  bool          `json:"active"`
	BookingsCount           int           `json:"bookings_count"`
	// BookingURL is the public hosted URL to share with humans.
	BookingURL string    `json:"booking_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateBookingPageParams are the parameters for creating a booking page.
type CreateBookingPageParams struct {
	CalendarID              string        `json:"calendar_id"`
	Title                   string        `json:"title"`
	Description             string        `json:"description,omitempty"`
	DurationMinutes         int           `json:"duration_minutes,omitempty"`
	BufferMinutes           int           `json:"buffer_minutes,omitempty"`
	WindowDays              int           `json:"window_days,omitempty"`
	MinNoticeMinutes        int           `json:"min_notice_minutes,omitempty"`
	Timezone                string        `json:"timezone,omitempty"`
	AvailabilityConstraints *WorkingHours `json:"availability_constraints,omitempty"`
	Active                  *bool         `json:"active,omitempty"`
}

// UpdateBookingPageParams are the parameters for updating a booking page.
type UpdateBookingPageParams struct {
	Title                   *string       `json:"title,omitempty"`
	Description             *string       `json:"description,omitempty"`
	DurationMinutes         *int          `json:"duration_minutes,omitempty"`
	BufferMinutes           *int          `json:"buffer_minutes,omitempty"`
	WindowDays              *int          `json:"window_days,omitempty"`
	MinNoticeMinutes        *int          `json:"min_notice_minutes,omitempty"`
	Timezone                *string       `json:"timezone,omitempty"`
	AvailabilityConstraints *WorkingHours `json:"availability_constraints,omitempty"`
	Active                  *bool         `json:"active,omitempty"`
}

// ListBookingPagesParams are the query parameters for listing booking pages.
type ListBookingPagesParams struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
	// CalendarID lists only the pages bound to this calendar.
	CalendarID *string `url:"calendar_id,omitempty"`
}

// Booking is a booking made through a booking page: the event it created plus
// the invitee's details from the booking form.
type Booking struct {
	// ID is the event the booking created.
	ID            string    `json:"id"`
	BookingPageID string    `json:"booking_page_id"`
	CalendarID    string    `json:"calendar_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	// Status is the event status, e.g. "confirmed", or "cancelled" once cancelled.
	Status       string  `json:"status"`
	InviteeName  *string `json:"invitee_name"`
	InviteeEmail *string `json:"invitee_email"`
	InviteeNotes *string `json:"invitee_notes"`
	// CreatedAt is when the booking was made.
	CreatedAt time.Time `json:"created_at"`
}

// ListBookingsParams are the query parameters for listing a booking page's bookings.
type ListBookingsParams struct {
	// Limit is the page size (1-100, default 50).
	Limit int `url:"limit,omitempty"`
	// Cursor is the NextCursor from the previous page.
	Cursor *string `url:"cursor,omitempty"`
}

// BookingList is one page of bookings.
type BookingList struct {
	Data []Booking `json:"data"`
	// NextCursor is passed as ListBookingsParams.Cursor for the next page; nil on the last page.
	NextCursor *string `json:"next_cursor"`
}

// --- Event ---

// EventStatus is the status of a calendar event.
type EventStatus string

const (
	EventStatusConfirmed EventStatus = "confirmed"
	EventStatusTentative EventStatus = "tentative"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusHold      EventStatus = "hold"
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
	CalendarID  string      `json:"calendarId"`
	OrgID       string      `json:"orgId"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     time.Time   `json:"endTime"`
	AllDay      bool        `json:"allDay"`
	Status      EventStatus `json:"status"`
	Source      EventSource `json:"source"`
	// Reminders are offsets in minutes before start_time at which an
	// event.reminder webhook fires and a VALARM is emitted in the iCal feed
	// (e.g. []int{10, 1440}). May be null/absent when the event inherits the
	// calendar default; an empty slice means no reminders.
	Reminders []int                  `json:"reminders"`
	Metadata  map[string]interface{} `json:"metadata"`
	// HoldExpiresAt and HoldPriority are populated only for temporal-hold
	// events (status "hold"); nil for regular events.
	HoldExpiresAt *time.Time `json:"holdExpiresAt"`
	HoldPriority  *int       `json:"holdPriority"`
	// RecurrenceRule is the RFC 5545 RRULE subset (no "RRULE:" prefix) when
	// this event is a recurring series master, e.g.
	// "FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12". nil for one-off events.
	RecurrenceRule *string `json:"recurrenceRule"`
	// RecurrenceExdates are the ISO 8601 starts of individually cancelled
	// occurrences (EXDATE) on a recurring series master.
	RecurrenceExdates []string `json:"recurrenceExdates"`
	// RecurringEventID is present only on expanded instances (expand=true):
	// the id of the recurring series master this instance belongs to.
	RecurringEventID string `json:"recurringEventId,omitempty"`
	// OriginalStartTime is present only on expanded instances: the occurrence
	// start this instance was generated for.
	OriginalStartTime string    `json:"originalStartTime,omitempty"`
	DeletedAt         *string   `json:"deletedAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
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
	Reminders []int `json:"reminders,omitempty"`
	// RecurrenceRule makes the event a recurring series master. RFC 5545 RRULE
	// subset without the "RRULE:" prefix, e.g. "FREQ=WEEKLY;BYDAY=MO,WE;COUNT=12".
	// Not allowed when Status is "hold".
	RecurrenceRule *string                `json:"recurrence_rule,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
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
	Reminders []int `json:"reminders,omitempty"`
	// RecurrenceRule sets or changes the series rule — a full-series edit that
	// re-anchors the series (RFC 5545 RRULE subset, no "RRULE:" prefix). The
	// REST API also accepts an explicit JSON null to clear the rule and make
	// the event one-off; like other nullable update fields this typed struct
	// only expresses set/change (omitted when nil).
	RecurrenceRule *string                `json:"recurrence_rule,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ListEventsParams are the query parameters for listing events.
type ListEventsParams struct {
	CalendarID  *string      `url:"calendar_id,omitempty"`
	AgentID     *string      `url:"agent_id,omitempty"`
	StartAfter  *string      `url:"start_after,omitempty"`
	StartBefore *string      `url:"start_before,omitempty"`
	Status      *EventStatus `url:"status,omitempty"`
	Source      *EventSource `url:"source,omitempty"`
	// Expand expands recurring series masters into concrete instances within
	// the query window. When true the API requires both StartAfter and
	// StartBefore, at most 366 days apart. Expanded instances carry
	// RecurringEventID and OriginalStartTime.
	Expand *bool `url:"expand,omitempty"`
	Limit  int   `url:"limit,omitempty"`
	Offset int   `url:"offset,omitempty"`
}

// ListOrgEventsParams filters Events.ListAll: events across all of the
// organization's calendars (GET /v1/events).
type ListOrgEventsParams struct {
	// CalendarIDs narrows to these calendars (at most 50), sent as repeated
	// calendar_id parameters.
	CalendarIDs []string     `url:"-"`
	AgentID     *string      `url:"agent_id,omitempty"`
	StartAfter  *string      `url:"start_after,omitempty"`
	StartBefore *string      `url:"start_before,omitempty"`
	Status      *EventStatus `url:"status,omitempty"`
	Source      *EventSource `url:"source,omitempty"`
	// Q searches titles case-insensitively. It requires StartAfter and
	// StartBefore at most 31 days apart.
	Q      *string `url:"q,omitempty"`
	Expand *bool   `url:"expand,omitempty"`
	// Limit is at most 100.
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
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
	Start string `json:"start"`
	End   string `json:"end"`
}

type AvailabilityBusyBlock struct {
	Start      string `json:"start"`
	End        string `json:"end"`
	CalendarID string `json:"calendar_id,omitempty"`
}

type AvailabilitySources struct {
	Chronary     string  `json:"chronary"`
	External     string  `json:"external"`
	LastSyncedAt *string `json:"last_synced_at"`
}

type AvailabilityWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AvailabilityResponse is the response from availability endpoints.
type AvailabilityResponse struct {
	Agents            []string                           `json:"agents,omitempty"`
	CalendarID        string                             `json:"calendar_id,omitempty"`
	Slots             []AvailabilitySlot                 `json:"slots"`
	PerAgentBusy      map[string][]AvailabilityBusyBlock `json:"per_agent_busy,omitempty"`
	Busy              []AvailabilityBusyBlock            `json:"busy,omitempty"`
	AvailabilityState string                             `json:"availability_state"`
	Sources           AvailabilitySources                `json:"sources"`
	Warnings          []AvailabilityWarning              `json:"warnings"`
}

// AvailabilityParams are the parameters for single-resource availability queries.
type AvailabilityParams struct {
	Start        string        `url:"start"`
	End          string        `url:"end"`
	SlotDuration *SlotDuration `url:"slot_duration,omitempty"`
	IncludeBusy  *bool         `url:"include_busy,omitempty"`
	AllowStale   *bool         `url:"allow_stale,omitempty"`
}

// CrossAgentAvailabilityParams are the parameters for cross-agent availability queries.
type CrossAgentAvailabilityParams struct {
	Agents       []string      `url:"agents"`
	Start        string        `url:"start"`
	End          string        `url:"end"`
	SlotDuration *SlotDuration `url:"slot_duration,omitempty"`
	Calendars    []string      `url:"calendars,omitempty"`
	IncludeBusy  *bool         `url:"include_busy,omitempty"`
	AllowStale   *bool         `url:"allow_stale,omitempty"`
}

// --- Human calendar connection setup links ---

type CreateConnectionLinkParams struct {
	Capabilities      []string `json:"capabilities"`
	PublicationPolicy string   `json:"publication_policy,omitempty"`
}

type ConnectionLink struct {
	ID           string  `json:"id"`
	CalendarID   string  `json:"calendar_id"`
	SetupURL     *string `json:"setup_url,omitempty"`
	Status       string  `json:"status"`
	ExpiresAt    string  `json:"expires_at"`
	ConnectionID *string `json:"connection_id,omitempty"`
	Reused       bool    `json:"reused,omitempty"`
	// Setup details (#1196).
	RequestingAgentID          *string  `json:"requesting_agent_id,omitempty"`
	RequestedCapabilities      []string `json:"requested_capabilities,omitempty"`
	RequestedPublicationPolicy *string  `json:"requested_publication_policy,omitempty"`
	CreatedAt                  string   `json:"created_at,omitempty"`
}

// ListConnectionLinksParams filters ConnectionLinks.List.
type ListConnectionLinksParams struct {
	// Status is awaiting_human, in_progress, completed, declined, expired or cancelled.
	Status     string
	CalendarID string
	// Limit is 1–100; the API defaults to 50.
	Limit int
}

// ConnectionLinkList is the ConnectionLinks.List response.
type ConnectionLinkList struct {
	Data []ConnectionLink `json:"data"`
}

// --- Webhook ---

// Webhook represents a Chronary webhook subscription.
type Webhook struct {
	ID                  string   `json:"id"`
	OrgID               string   `json:"orgId"`
	URL                 string   `json:"url"`
	Events              []string `json:"events"`
	Active              bool     `json:"active"`
	ConsecutiveFailures int      `json:"consecutiveFailures"`
	FirstFailureAt      *string  `json:"firstFailureAt"`
	// DeletionPausedAt is set while a scheduled account deletion has paused
	// this endpoint (Active is then false).
	DeletionPausedAt *string `json:"deletionPausedAt,omitempty"`
	// DeactivatedAt is when the endpoint last went inactive; nil while it is
	// active. Cleared when it is re-enabled.
	DeactivatedAt *string `json:"deactivatedAt,omitempty"`
	// DeactivatedReason is why the endpoint is inactive (one of the
	// WebhookDeactivated* constants); nil while it is active.
	DeactivatedReason *string `json:"deactivatedReason,omitempty"`
	// PreviousSecretExpiresAt is set while a secret rotation's grace window is
	// open: when the replaced secret stops working. Until then deliveries also
	// carry X-Signature-Previous. Nil when no window is open.
	PreviousSecretExpiresAt *string   `json:"previousSecretExpiresAt,omitempty"`
	CreatedAt               time.Time `json:"createdAt"`
}

// Reasons a webhook endpoint is inactive (Webhook.DeactivatedReason).
const (
	WebhookDeactivatedManual               = "manual"
	WebhookDeactivatedCircuitBreaker       = "circuit_breaker"
	WebhookDeactivatedUnsafeDestinationURL = "unsafe_destination_url"
	WebhookDeactivatedAccountDeletion      = "account_deletion"
)

// SendTestWebhookParams selects the sample payload Webhooks.SendTest sends.
// The endpoint must subscribe to EventType.
type SendTestWebhookParams struct {
	EventType string `json:"event_type"`
}

// WebhookTestSent is the queued test delivery Webhooks.SendTest returns.
type WebhookTestSent struct {
	DeliveryID string `json:"delivery_id"`
	WebhookID  string `json:"webhook_id"`
	EventType  string `json:"event_type"`
	Status     string `json:"status"`
	IsTest     bool   `json:"is_test"`
}

// WebhookReplaySent is the new delivery Webhooks.ReplayDelivery queued.
type WebhookReplaySent struct {
	// DeliveryID is the new delivery; poll GetDelivery with it.
	DeliveryID string `json:"delivery_id"`
	WebhookID  string `json:"webhook_id"`
	EventType  string `json:"event_type"`
	Status     string `json:"status"`
	// ReplayOfDeliveryID is the first delivery of the chain, sent as the
	// X-Chronary-Replay-Of header.
	ReplayOfDeliveryID string `json:"replay_of_delivery_id"`
}

// When the replaced secret stops working after Webhooks.RotateSecret
// (RotateWebhookSecretParams.ExpireOldIn).
const (
	WebhookSecretExpireNow = "now"
	WebhookSecretExpire1h  = "1h"
	WebhookSecretExpire24h = "24h"
	WebhookSecretExpire7d  = "7d"
)

// RotateWebhookSecretParams configures Webhooks.RotateSecret. An empty
// ExpireOldIn uses the API default, 24h.
type RotateWebhookSecretParams struct {
	ExpireOldIn string `json:"expire_old_in,omitempty"`
}

// RotatedWebhookSecret is the new signing secret Webhooks.RotateSecret returns
// once. PreviousSecretExpiresAt is nil when ExpireOldIn was "now".
type RotatedWebhookSecret struct {
	WebhookID               string  `json:"webhook_id"`
	Secret                  string  `json:"secret"`
	PreviousSecretExpiresAt *string `json:"previous_secret_expires_at"`
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
// period. Informational — not gated by any plan limit. Created, Confirmed
// and Expired are monthly counters; Expired covers all three end-of-hold
// paths: TTL expiry, manual release, and priority-bump. Active is a live
// count of unexpired holds, not a monthly counter.
type HoldsUsage struct {
	Created   int `json:"created"`
	Confirmed int `json:"confirmed"`
	Expired   int `json:"expired"`
	Active    int `json:"active"`
}

// CrossCalendarQueriesUsage counts availability requests that touched more
// than one calendar in the current period. Informational — gated separately
// by the cross_calendar_availability capability, not by this counter.
type CrossCalendarQueriesUsage struct {
	Used int `json:"used"`
}

// ScopedKeysUsage counts active agent-scoped keys for the current period. Its
// server shape uses `count` (not `used`, for backwards-compat with the #308
// informational shape) and a nullable `limit` (null = unlimited on Scale).
type ScopedKeysUsage struct {
	Count int  `json:"count"`
	Limit *int `json:"limit"`
}

// Usage is the response from the usage endpoint.
type Usage struct {
	PeriodStart          string                    `json:"period_start"`
	PeriodEnd            string                    `json:"period_end"`
	ResetsAt             string                    `json:"resets_at"`
	GeneratedAt          string                    `json:"generated_at"`
	Plan                 string                    `json:"plan"`
	Agents               UsageCounter              `json:"agents"`
	Calendars            UsageCounter              `json:"calendars"`
	Events               UsageCounter              `json:"events"`
	APICalls             UsageCounter              `json:"api_calls"`
	Webhooks             UsageCounter              `json:"webhooks"`
	WebhookEndpoints     UsageCounter              `json:"webhook_endpoints"`
	AvailabilityQueries  UsageCounter              `json:"availability_queries"`
	ICalSubscriptions    UsageCounter              `json:"ical_subscriptions"`
	Proposals            UsageCounter              `json:"proposals"`
	RecurringEvents      UsageCounter              `json:"recurring_events"`
	BookingPages         UsageCounter              `json:"booking_pages"`
	ScopedKeys           ScopedKeysUsage           `json:"scoped_keys"`
	Holds                HoldsUsage                `json:"holds"`
	CrossCalendarQueries CrossCalendarQueriesUsage `json:"cross_calendar_queries"`
}

// UsageHistoryMetrics holds the monthly counter totals for one UTC calendar month.
type UsageHistoryMetrics struct {
	APICalls             int `json:"api_calls"`
	Events               int `json:"events"`
	WebhookDeliveries    int `json:"webhook_deliveries"`
	AvailabilityQueries  int `json:"availability_queries"`
	Proposals            int `json:"proposals"`
	HoldsCreated         int `json:"holds_created"`
	HoldsConfirmed       int `json:"holds_confirmed"`
	HoldsExpired         int `json:"holds_expired"`
	CrossCalendarQueries int `json:"cross_calendar_queries"`
}

// UsageHistoryMonth is one UTC calendar month of usage history.
type UsageHistoryMonth struct {
	// Period is the UTC calendar month, YYYY-MM.
	Period      string `json:"period"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
	// IsCurrent is true only for the current, partial month (the last entry).
	IsCurrent bool                `json:"is_current"`
	Metrics   UsageHistoryMetrics `json:"metrics"`
}

// UsageHistoryLimits holds the current plan's limits for the limited monthly
// metrics. A nil limit means unlimited. Historical limits are not recorded.
type UsageHistoryLimits struct {
	APICalls            *int `json:"api_calls"`
	Events              *int `json:"events"`
	WebhookDeliveries   *int `json:"webhook_deliveries"`
	AvailabilityQueries *int `json:"availability_queries"`
	Proposals           *int `json:"proposals"`
}

// UsageHistory is the response from the usage history endpoint: monthly
// totals, oldest month first, ending with the current month.
type UsageHistory struct {
	Plan   string              `json:"plan"`
	Months int                 `json:"months"`
	Limits UsageHistoryLimits  `json:"limits"`
	Data   []UsageHistoryMonth `json:"data"`
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
	WebhookEventOrgPlanChanged     WebhookEventType = "org.plan_changed"
	// WebhookEventCalendarConnectionReauthorizationRequired fires once when a
	// connected human calendar stops accepting access (#1196).
	WebhookEventCalendarConnectionReauthorizationRequired WebhookEventType = "calendar_connection.reauthorization_required"
	// WebhookEventUsageThresholdCrossed fires once per metric, period and
	// threshold (50, 80, 90, 95 or 100 percent) of a monthly usage limit (#1211).
	WebhookEventUsageThresholdCrossed WebhookEventType = "usage.threshold_crossed"
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
	WebhookEventWebhookDeactivated, WebhookEventOrgPlanChanged,
	WebhookEventCalendarConnectionReauthorizationRequired,
	WebhookEventUsageThresholdCrossed,
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
	// WebhookDeliveryDropped is a delivery the Free plan monthly limit
	// stopped: never sent, counted or retried; LastErrorCode is
	// "quota_exceeded".
	WebhookDeliveryDropped WebhookDeliveryStatus = "dropped"
)

// WebhookDelivery is a single delivery attempt for a webhook subscription.
type WebhookDelivery struct {
	ID             string                `json:"id"`
	SubscriptionID string                `json:"subscription_id"`
	EventType      string                `json:"event_type"`
	Status         WebhookDeliveryStatus `json:"status"`
	Attempts       int                   `json:"attempts"`
	LastAttemptAt  *string               `json:"last_attempt_at"`
	NextRetryAt    *string               `json:"next_retry_at"`
	CreatedAt      string                `json:"created_at"`
	// LastResponseStatus is the HTTP status of the most recent attempt; nil
	// when no response was received.
	LastResponseStatus *int `json:"last_response_status"`
	// LastErrorCode is why the most recent attempt failed; nil on success.
	LastErrorCode *string `json:"last_error_code"`
	// LastDurationMs is the most recent attempt's time to response headers.
	LastDurationMs *int `json:"last_duration_ms"`
	// IsTest is true for a test event sent with Webhooks.SendTest.
	IsTest bool `json:"is_test"`
	// ReplayOfDeliveryID is set on a replay sent with Webhooks.ReplayDelivery:
	// the first delivery of the chain, the original even when a replay is
	// replayed. Nil otherwise.
	ReplayOfDeliveryID *string                `json:"replay_of_delivery_id"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
}

// Webhook attempt error codes. A delivery's LastErrorCode may also be
// "endpoint_disabled" or "quota_exceeded".
const (
	WebhookAttemptErrorTimeout              = "timeout"
	WebhookAttemptErrorConnection           = "connection_error"
	WebhookAttemptErrorTLS                  = "tls_error"
	WebhookAttemptErrorRedirectNotFollowed  = "redirect_not_followed"
	WebhookAttemptErrorUnsafeDestinationURL = "unsafe_destination_url"
	WebhookAttemptErrorNon2xx               = "non_2xx"
)

// WebhookDeliveryAttempt is one HTTP attempt of a webhook delivery.
type WebhookDeliveryAttempt struct {
	ID             string  `json:"id"`
	AttemptNumber  int     `json:"attempt_number"`
	StartedAt      string  `json:"started_at"`
	DurationMs     *int    `json:"duration_ms"`
	ResponseStatus *int    `json:"response_status"`
	ResponseBody   *string `json:"response_body"`
	ErrorCode      *string `json:"error_code"`
}

// WebhookDeliveryDetail is one webhook delivery with its payload and every
// attempt, oldest first.
type WebhookDeliveryDetail struct {
	ID                 string                   `json:"id"`
	SubscriptionID     string                   `json:"subscription_id"`
	EventType          string                   `json:"event_type"`
	Status             WebhookDeliveryStatus    `json:"status"`
	Attempts           int                      `json:"attempts"`
	MaxAttempts        int                      `json:"max_attempts"`
	LastAttemptAt      *string                  `json:"last_attempt_at"`
	NextRetryAt        *string                  `json:"next_retry_at"`
	CreatedAt          string                   `json:"created_at"`
	LastResponseStatus *int                     `json:"last_response_status"`
	LastErrorCode      *string                  `json:"last_error_code"`
	LastDurationMs     *int                     `json:"last_duration_ms"`
	IsTest             bool                     `json:"is_test"`
	ReplayOfDeliveryID *string                  `json:"replay_of_delivery_id"`
	Payload            interface{}              `json:"payload"`
	PayloadRedacted    bool                     `json:"payload_redacted"`
	RequestHeaders     map[string]string        `json:"request_headers"`
	AttemptLog         []WebhookDeliveryAttempt `json:"attempt_log"`
}

// WebhookDeliveryStats is the rollup returned alongside a deliveries list.
type WebhookDeliveryStats struct {
	Pending   int `json:"pending"`
	Delivered int `json:"delivered"`
	Failed    int `json:"failed"`
	// Dropped counts deliveries the Free plan monthly limit stopped.
	Dropped int `json:"dropped"`
}

// WebhookDeliveryListResponse is the response from listing webhook deliveries.
type WebhookDeliveryListResponse struct {
	Data []WebhookDelivery `json:"data"`
	// Total counts deliveries matching the filters, ignoring Cursor and Offset.
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	// NextCursor is passed as Cursor for the next page; nil on the last page.
	NextCursor *string              `json:"next_cursor"`
	HasMore    bool                 `json:"has_more"`
	Stats      WebhookDeliveryStats `json:"stats"`
}

// ListWebhookDeliveriesParams are query parameters for listing webhook deliveries.
type ListWebhookDeliveriesParams struct {
	Status         *WebhookDeliveryStatus `url:"status,omitempty"`
	IncludePayload *bool                  `url:"include_payload,omitempty"`
	Limit          int                    `url:"limit,omitempty"`
	// Offset cannot be combined with Cursor.
	Offset int `url:"offset,omitempty"`
	// EventType is an exact type ("event.created") or a prefix ending in a
	// dot ("event.").
	EventType *string `url:"event_type,omitempty"`
	// CreatedAfter (inclusive) and CreatedBefore (exclusive) are ISO 8601 UTC.
	CreatedAfter  *string `url:"created_after,omitempty"`
	CreatedBefore *string `url:"created_before,omitempty"`
	DeliveryID    *string `url:"delivery_id,omitempty"`
	// IsTest true lists only test events; false leaves them out.
	IsTest *bool `url:"is_test,omitempty"`
	// Cursor is the NextCursor of the previous page.
	Cursor *string `url:"cursor,omitempty"`
}

// --- Scheduling proposals ---

// ProposalStatus is the lifecycle state of a scheduling proposal.
type ProposalStatus string

const (
	ProposalStatusPending ProposalStatus = "pending"
	// ProposalStatusResolving is the brief state while a resolve creates the
	// event; it returns to pending if the resolve fails.
	ProposalStatusResolving ProposalStatus = "resolving"
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
	// LastUsedAt is approximate (updated at most every 5 minutes); nil if never used.
	LastUsedAt        *string `json:"last_used_at"`
	LastUsedUserAgent *string `json:"last_used_user_agent"`
	// CreatedVia is one of signup, console, api, mcp or unknown.
	CreatedVia string `json:"created_via"`
	// CreatedByKeyPrefix is the creating Organization key's prefix (api and mcp).
	CreatedByKeyPrefix *string `json:"created_by_key_prefix"`
	// Note is a free-text note about the key, set in the console.
	Note *string `json:"note"`
	// ExpiresAt is when the key stops working (401 key_expired); nil = never expires.
	ExpiresAt *string `json:"expires_at"`
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
	// ExpiresIn is one of 30d, 60d, 90d, 180d or 365d; nil = the key never expires.
	ExpiresIn *string `json:"expires_in,omitempty"`
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
	WebhookEndpoints    *int `json:"webhook_endpoints"`
	ScopedKeys          *int `json:"scoped_keys"`
}

// Plan is a single tier in the public plan catalog.
type Plan struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Tagline         string      `json:"tagline"`
	Price           *int        `json:"price"`
	Currency        *string     `json:"currency"`
	AnnualPrice     *int        `json:"annual_price,omitempty"`
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
	Verified bool `json:"verified"`
	// SandboxPromotion is "promoted" when the org left the Sandbox tier,
	// "capacity_full" when the email is verified but early access is full (the
	// org keeps Sandbox limits until a spot opens), or nil when the org was not
	// on Sandbox.
	SandboxPromotion *string `json:"sandbox_promotion"`
	Message          string  `json:"message"`
}

// --- Feedback ---

// FeedbackType categorizes a feedback submission.
type FeedbackType string

const (
	FeedbackTypeBug      FeedbackType = "bug"
	FeedbackTypeFeature  FeedbackType = "feature"
	FeedbackTypeFriction FeedbackType = "friction"
	FeedbackTypeSupport  FeedbackType = "support"
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
	// Nil for console sessions and the system.
	ActorKeyPrefix *string `json:"actor_key_prefix"`
	// ActorType is who acted: "session", "org_key", "agent_key" or "system".
	// Nil on older entries that could not be attributed.
	ActorType *string `json:"actor_type"`
	// ActorID is the API key id for a key, the organization id for a console session.
	ActorID *string `json:"actor_id"`
	// ActorLabel is the key label (or key prefix), the session email, or "Chronary".
	ActorLabel *string `json:"actor_label"`
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
	// RetentionDays is the plan's audit-log retention window (days). Nil = no fixed window (Custom: per contract).
	RetentionDays *int `json:"retention_days"`
	// RangeClamped is true when the requested from-date was outside the retention window.
	RangeClamped bool `json:"range_clamped"`
}

// --- Billing ---

// BillingPrice is the org's current subscription price. Amount is in minor
// units (2900 = $29).
type BillingPrice struct {
	Amount       int    `json:"amount"`
	Currency     string `json:"currency"`
	Interval     string `json:"interval"`
	TaxExclusive bool   `json:"tax_exclusive"`
}

// BillingUpgradeOption is a plan/interval a human can upgrade to.
type BillingUpgradeOption struct {
	Plan             string `json:"plan"`
	Interval         string `json:"interval"`
	Amount           int    `json:"amount"`
	Currency         string `json:"currency"`
	SavingsVsMonthly int    `json:"savings_vs_monthly"`
	Via              string `json:"via,omitempty"`
}

// UpgradeRequest is an agent-initiated request for a human to upgrade the
// account. Status is one of pending, completed, declined, canceled, expired.
type UpgradeRequest struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Plan     string `json:"plan"`
	Interval string `json:"interval"`
	// Reason is agent free text; render it as escaped plain text only.
	Reason  *string `json:"reason"`
	Feature *string `json:"feature"`
	Metric  *string `json:"metric"`
	// RequestedVia is "api" or "mcp".
	RequestedVia string `json:"requested_via"`
	// RequestedByKeyPrefix is the first 20 chars of the requesting key.
	RequestedByKeyPrefix *string `json:"requested_by_key_prefix"`
	// ApprovalURL is the link to give your human. Nil once the request is terminal.
	ApprovalURL *string `json:"approval_url"`
	CreatedAt   string  `json:"created_at"`
	ExpiresAt   string  `json:"expires_at"`
	CompletedAt *string `json:"completed_at"`
	DeclinedAt  *string `json:"declined_at"`
	CanceledAt  *string `json:"canceled_at"`
	// OwnerNotification and Reused are set on create responses only.
	OwnerNotification string `json:"owner_notification,omitempty"`
	Reused            bool   `json:"reused,omitempty"`
}

// BillingSummary is the read-only billing state of the org.
type BillingSummary struct {
	Plan                string                 `json:"plan"`
	PlanSource          *string                `json:"plan_source"`
	PlanExpiresAt       *string                `json:"plan_expires_at"`
	Status              *string                `json:"status"`
	Interval            *string                `json:"interval"`
	CurrentPeriodEnd    *string                `json:"current_period_end"`
	CancelAt            *string                `json:"cancel_at"`
	CancelAtPeriodEnd   bool                   `json:"cancel_at_period_end"`
	RenewsAt            *string                `json:"renews_at"`
	PendingPlan         *string                `json:"pending_plan,omitempty"`
	PendingPlanInterval *string                `json:"pending_plan_interval,omitempty"`
	PendingPlanAt       *string                `json:"pending_plan_at,omitempty"`
	Price               *BillingPrice          `json:"price"`
	HasBillingAccount   bool                   `json:"has_billing_account"`
	SelfServeUpgrade    bool                   `json:"self_serve_upgrade"`
	UpgradeOptions      []BillingUpgradeOption `json:"upgrade_options"`
	UpgradeURL          *string                `json:"upgrade_url"`
	ManageURL           string                 `json:"manage_url"`
	ContactURL          *string                `json:"contact_url"`
	// PendingUpgradeRequest is set only on Billing.Get.
	PendingUpgradeRequest *UpgradeRequest `json:"pending_upgrade_request"`
	Offer                 interface{}     `json:"offer"`
	// PastDue is set while Status is past_due: the failed payment.
	PastDue *BillingPastDue `json:"past_due"`
	// Credit is account credit, applied automatically to the next invoices while
	// the subscription renews; nil when there is none.
	Credit *BillingCredit `json:"credit"`
}

// BillingCredit is account credit, for example the unused part of a yearly
// period after a switch to monthly billing.
type BillingCredit struct {
	// Amount is in minor units (cents), always positive.
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// BillingPastDue is the failed payment behind a past_due subscription.
type BillingPastDue struct {
	// AmountDue is in minor units (cents), tax included.
	AmountDue int    `json:"amount_due"`
	Currency  string `json:"currency"`
	FailedAt  string `json:"failed_at"`
	// NextPaymentAttempt is nil when no automatic retries are left.
	NextPaymentAttempt *string `json:"next_payment_attempt"`
	// HostedInvoiceURL is where a human can pay the open invoice.
	HostedInvoiceURL *string `json:"hosted_invoice_url"`
	AttemptCount     int     `json:"attempt_count"`
}

// CreateUpgradeRequestParams are the inputs to Billing.CreateUpgradeRequest.
type CreateUpgradeRequestParams struct {
	// Interval is "month" (default) or "year".
	Interval string `json:"interval,omitempty"`
	// Reason is shown to the human approver (max 500 chars).
	Reason string `json:"reason,omitempty"`
	// Feature is the plan-gated feature that was blocked, e.g. "scoped_keys".
	Feature string `json:"feature,omitempty"`
	// Metric is the quota metric that was exhausted, e.g. "api_calls".
	Metric string `json:"metric,omitempty"`
}
