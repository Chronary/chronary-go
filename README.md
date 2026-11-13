# chronary-go

Official Go SDK for the [Chronary](https://chronary.ai) calendar-as-a-service API.

## Installation

```bash
go get github.com/Chronary/chronary-go
```

Requires Go 1.24 or later.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    chronary "github.com/Chronary/chronary-go"
)

func main() {
    client, err := chronary.NewClient(
        chronary.WithAPIKey("chr_sk_..."),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Register your agent with Chronary
    agent, err := client.Agents.Create(ctx, &chronary.CreateAgentParams{
        Name: "Meeting Bot",
        Type: chronary.AgentTypeAI,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a calendar
    cal, err := client.Calendars.CreateForAgent(ctx, agent.ID, &chronary.CreateCalendarParams{
        Name:     "Team Meetings",
        Timezone: "America/New_York",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create an event
    event, err := client.Events.Create(ctx, cal.ID, &chronary.CreateEventParams{
        Title:     "Standup",
        StartTime: "2026-04-15T09:00:00Z",
        EndTime:   "2026-04-15T09:30:00Z",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created event: %s\n", event.ID)
}
```

## Pagination

List methods return a `PageIterator` that supports Go's range-over-func:

```go
for agent, err := range client.Agents.List(ctx, nil).All(ctx) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(agent.Name)
}
```

Or collect all results into a slice:

```go
agents, err := client.Agents.List(ctx, nil).Collect(ctx)
```

Or fetch a specific page:

```go
page, err := client.Agents.List(ctx, &chronary.ListAgentsParams{
    Limit: 10,
}).GetPage(ctx, 0)
```

## Error Handling

Errors from the API are returned as `*chronary.Error` and support `errors.Is` and `errors.As`:

```go
_, err := client.Agents.Get(ctx, "agt_nonexistent")
if errors.Is(err, chronary.ErrNotFound) {
    fmt.Println("Agent not found")
}

var chronErr *chronary.Error
if errors.As(err, &chronErr) {
    fmt.Printf("Error type: %s, request ID: %s\n", chronErr.Type, chronErr.RequestID)
}
```

Sentinel errors: `ErrAuthentication`, `ErrNotFound`, `ErrValidation`, `ErrRateLimit`, `ErrQuotaExceeded`, `ErrTimeout`, `ErrConnection`.

## Reminders

Reminder fields on the create and update params are `*[]int` (minutes before the event starts), so "not set", "none" and "inherit" stay distinct:

```go
// 10 minutes and 1 day before
event, err := client.Events.Create(ctx, cal.ID, &chronary.CreateEventParams{
    Title:     "Standup",
    StartTime: "2026-04-15T09:00:00Z",
    EndTime:   "2026-04-15T09:30:00Z",
    Reminders: chronary.Reminders(10, 1440),
})

// No reminders on this event (sends [])
_, err = client.Events.Update(ctx, cal.ID, event.ID, &chronary.UpdateEventParams{Reminders: chronary.NoReminders()})

// Go back to the calendar's default (sends null)
_, err = client.Events.Update(ctx, cal.ID, event.ID, &chronary.UpdateEventParams{Reminders: chronary.InheritReminders()})

// Clear a calendar's default reminders
_, err = client.Calendars.Update(ctx, cal.ID, &chronary.UpdateCalendarParams{DefaultReminders: chronary.NoReminders()})
```

Leaving the field nil omits it: on a create the event inherits the calendar default (then the system default of 10 minutes); on an update the schedule is unchanged.

## Availability

```go
// Check a single agent's availability
avail, err := client.Availability.ForAgent(ctx, "agt_1", &chronary.AvailabilityParams{
    Start: "2026-04-15T09:00:00Z",
    End:   "2026-04-15T17:00:00Z",
})

// Check across multiple agents
avail, err := client.Availability.Check(ctx, &chronary.CrossAgentAvailabilityParams{
    Agents: []string{"agt_1", "agt_2"},
    Start:  "2026-04-15T09:00:00Z",
    End:    "2026-04-15T17:00:00Z",
})
```

## Webhook Verification

Verify webhook signatures without creating a client:

```go
err := chronary.VerifySignature(payload, req.Header, webhookSecret)
if err != nil {
    http.Error(w, "Invalid signature", 403)
    return
}

// Or verify and parse in one step
event, err := chronary.ConstructEvent(payload, req.Header, webhookSecret)
if err != nil {
    http.Error(w, "Invalid webhook", 403)
    return
}
fmt.Printf("Received: %s\n", event.Type)
```

## Configuration

```go
client, err := chronary.NewClient(
    chronary.WithAPIKey("chr_sk_..."),       // or set CHRONARY_API_KEY env var
    chronary.WithBaseURL("https://api.custom.ai"), // custom base URL
    chronary.WithTimeout(60 * time.Second),        // request timeout (default: 30s)
    chronary.WithMaxRetries(3),                    // retry count (default: 2)
    chronary.WithHTTPClient(customHTTPClient),     // custom http.Client
)
```

## Retries

The SDK automatically retries failed requests with exponential backoff for:
- HTTP 408 (Request Timeout)
- HTTP 429 (Rate Limited) -- respects `Retry-After` header
- HTTP 500, 502, 503, 504 (Server Errors)

Default: 2 retries. Disable with `WithMaxRetries(0)`.

## Resources

| Service | Methods |
|---------|---------|
| `Agents` | Create, Get, List, Update, Delete, ListCalendars, ListEvents, ListICalSubscriptions |
| `Calendars` | Create, CreateForAgent, Get, List, Update, Delete |
| `Events` | Create, Get, List, Update, Delete |
| `Availability` | ForAgent, ForCalendar, Check |
| `Webhooks` | Create, Get, List, Update, Delete |
| `ICalSubscriptions` | Create, Get, List, Update, Delete, Sync |
| `Usage` | Get |

## Security and support

Report vulnerabilities to [security@chronary.ai](mailto:security@chronary.ai) or through private vulnerability reporting on this package's GitHub repository. Please don't open a public issue. Read the full [security policy](https://github.com/Chronary/.github/blob/main/SECURITY.md).

Version 0.x is supported with security fixes until at least **April 2031**. Fixes ship in a new 0.x release, so keep up with the latest one. Supported versions and end dates for every Chronary package are listed in the [security policy](https://github.com/Chronary/.github/blob/main/SECURITY.md#supported-versions).

## License

Apache-2.0
