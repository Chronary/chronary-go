package chronary_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"testing"

	chronary "github.com/Chronary/chronary-go"
	"github.com/Chronary/chronary-go/internal/testutil"
)

// #2519: IDs used to be interpolated raw with fmt.Sprintf, so a caller-supplied
// ID like "../agents/agt_1" put a dot segment on the wire (which the edge and
// the Worker's URL parser resolve to /v1/agents/agt_1), "#" dropped the rest
// of the path and "?" added a query string. Every ID now goes through
// pathSegment: one escaped segment, and "", "." and ".." refused before any
// request is sent.

const escapeID = "../keys#x?y"
const escapeEncoded = "..%2Fkeys%23x%3Fy"

type pathCall struct {
	name   string
	call   func(c *chronary.Client, id string) error
	prefix string
}

func ignore[T any](_ T, err error) error { return err }

func pathCalls() []pathCall {
	ctx := context.Background()
	s := func(v string) *string { return &v }
	return []pathCall{
		{"Agents.Get", func(c *chronary.Client, id string) error { return ignore(c.Agents.Get(ctx, id)) }, "/v1/agents/"},
		{"Agents.Update", func(c *chronary.Client, id string) error {
			return ignore(c.Agents.Update(ctx, id, &chronary.UpdateAgentParams{}))
		}, "/v1/agents/"},
		{"Agents.Delete", func(c *chronary.Client, id string) error { return c.Agents.Delete(ctx, id) }, "/v1/agents/"},
		{"Agents.ListCalendars", func(c *chronary.Client, id string) error {
			return ignore(c.Agents.ListCalendars(ctx, id, nil).GetPage(ctx, 0))
		}, "/v1/agents/"},
		{"Agents.ListEvents", func(c *chronary.Client, id string) error {
			return ignore(c.Agents.ListEvents(ctx, id, nil).GetPage(ctx, 0))
		}, "/v1/agents/"},
		{"Agents.ListICalSubscriptions", func(c *chronary.Client, id string) error {
			return ignore(c.Agents.ListICalSubscriptions(ctx, id, nil).GetPage(ctx, 0))
		}, "/v1/agents/"},
		{"Availability.ForAgent", func(c *chronary.Client, id string) error {
			return ignore(c.Availability.ForAgent(ctx, id, &chronary.AvailabilityParams{}))
		}, "/v1/agents/"},
		{"Availability.ForCalendar", func(c *chronary.Client, id string) error {
			return ignore(c.Availability.ForCalendar(ctx, id, &chronary.AvailabilityParams{}))
		}, "/v1/calendars/"},
		{"Billing.GetUpgradeRequest", func(c *chronary.Client, id string) error { return ignore(c.Billing.GetUpgradeRequest(ctx, id)) }, "/v1/billing/upgrade-requests/"},
		{"Billing.CancelUpgradeRequest", func(c *chronary.Client, id string) error {
			return ignore(c.Billing.CancelUpgradeRequest(ctx, id))
		}, "/v1/billing/upgrade-requests/"},
		{"BookingPages.Get", func(c *chronary.Client, id string) error { return ignore(c.BookingPages.Get(ctx, id)) }, "/v1/booking-pages/"},
		{"BookingPages.Update", func(c *chronary.Client, id string) error {
			return ignore(c.BookingPages.Update(ctx, id, &chronary.UpdateBookingPageParams{}))
		}, "/v1/booking-pages/"},
		{"BookingPages.Delete", func(c *chronary.Client, id string) error { return c.BookingPages.Delete(ctx, id) }, "/v1/booking-pages/"},
		{"BookingPages.ListBookings", func(c *chronary.Client, id string) error {
			return ignore(c.BookingPages.ListBookings(ctx, id, nil))
		}, "/v1/booking-pages/"},
		{"BookingPages.RotateSlug", func(c *chronary.Client, id string) error { return ignore(c.BookingPages.RotateSlug(ctx, id)) }, "/v1/booking-pages/"},
		{"Calendars.CreateForAgent", func(c *chronary.Client, id string) error {
			return ignore(c.Calendars.CreateForAgent(ctx, id, &chronary.CreateCalendarParams{}))
		}, "/v1/agents/"},
		{"Calendars.Get", func(c *chronary.Client, id string) error { return ignore(c.Calendars.Get(ctx, id)) }, "/v1/calendars/"},
		{"Calendars.RotateICalToken", func(c *chronary.Client, id string) error { return ignore(c.Calendars.RotateICalToken(ctx, id)) }, "/v1/calendars/"},
		{"Calendars.Update", func(c *chronary.Client, id string) error {
			return ignore(c.Calendars.Update(ctx, id, &chronary.UpdateCalendarParams{}))
		}, "/v1/calendars/"},
		{"Calendars.Delete", func(c *chronary.Client, id string) error { return c.Calendars.Delete(ctx, id) }, "/v1/calendars/"},
		{"Calendars.GetContext", func(c *chronary.Client, id string) error { return ignore(c.Calendars.GetContext(ctx, id)) }, "/v1/calendars/"},
		{"Calendars.SetAvailabilityRules", func(c *chronary.Client, id string) error {
			return ignore(c.Calendars.SetAvailabilityRules(ctx, id, &chronary.SetAvailabilityRulesParams{}))
		}, "/v1/calendars/"},
		{"Calendars.GetAvailabilityRules", func(c *chronary.Client, id string) error {
			return ignore(c.Calendars.GetAvailabilityRules(ctx, id))
		}, "/v1/calendars/"},
		{"Calendars.DeleteAvailabilityRules", func(c *chronary.Client, id string) error {
			return c.Calendars.DeleteAvailabilityRules(ctx, id)
		}, "/v1/calendars/"},
		{"ConnectionLinks.Create", func(c *chronary.Client, id string) error {
			return ignore(c.ConnectionLinks.Create(ctx, id, &chronary.CreateConnectionLinkParams{}))
		}, "/v1/calendars/"},
		{"ConnectionLinks.Get", func(c *chronary.Client, id string) error { return ignore(c.ConnectionLinks.Get(ctx, id)) }, "/v1/connection-links/"},
		{"ConnectionLinks.Cancel", func(c *chronary.Client, id string) error { return c.ConnectionLinks.Cancel(ctx, id) }, "/v1/connection-links/"},
		{"Events.Create", func(c *chronary.Client, id string) error {
			return ignore(c.Events.Create(ctx, id, &chronary.CreateEventParams{}))
		}, "/v1/calendars/"},
		{"Events.Get(calendarID)", func(c *chronary.Client, id string) error { return ignore(c.Events.Get(ctx, id, "evt_1")) }, "/v1/calendars/"},
		{"Events.Get(eventID)", func(c *chronary.Client, id string) error { return ignore(c.Events.Get(ctx, "cal_1", id)) }, "/v1/calendars/cal_1/events/"},
		{"Events.List(calendarID)", func(c *chronary.Client, id string) error {
			return ignore(c.Events.List(ctx, &chronary.ListEventsParams{CalendarID: s(id)}).GetPage(ctx, 0))
		}, "/v1/calendars/"},
		{"Events.List(agentID)", func(c *chronary.Client, id string) error {
			return ignore(c.Events.List(ctx, &chronary.ListEventsParams{AgentID: s(id)}).GetPage(ctx, 0))
		}, "/v1/agents/"},
		{"Events.Update(eventID)", func(c *chronary.Client, id string) error {
			return ignore(c.Events.Update(ctx, "cal_1", id, &chronary.UpdateEventParams{}))
		}, "/v1/calendars/cal_1/events/"},
		{"Events.Delete(eventID)", func(c *chronary.Client, id string) error { return c.Events.Delete(ctx, "cal_1", id) }, "/v1/calendars/cal_1/events/"},
		{"Events.CancelOccurrence(eventID)", func(c *chronary.Client, id string) error {
			return ignore(c.Events.CancelOccurrence(ctx, "cal_1", id, "2026-01-01T00:00:00Z"))
		}, "/v1/calendars/cal_1/events/"},
		{"Events.GetByID", func(c *chronary.Client, id string) error { return ignore(c.Events.GetByID(ctx, id)) }, "/v1/events/"},
		{"Events.UpdateByID", func(c *chronary.Client, id string) error {
			return ignore(c.Events.UpdateByID(ctx, id, &chronary.UpdateEventParams{}))
		}, "/v1/events/"},
		{"Events.DeleteByID", func(c *chronary.Client, id string) error { return c.Events.DeleteByID(ctx, id) }, "/v1/events/"},
		{"Events.CancelOccurrenceByID", func(c *chronary.Client, id string) error {
			return ignore(c.Events.CancelOccurrenceByID(ctx, id, "2026-01-01T00:00:00Z"))
		}, "/v1/events/"},
		{"Events.Confirm", func(c *chronary.Client, id string) error { return ignore(c.Events.Confirm(ctx, id)) }, "/v1/events/"},
		{"Events.Release", func(c *chronary.Client, id string) error { return ignore(c.Events.Release(ctx, id)) }, "/v1/events/"},
		{"ICalSubscriptions.Create", func(c *chronary.Client, id string) error {
			return ignore(c.ICalSubscriptions.Create(ctx, id, &chronary.CreateICalSubscriptionParams{}))
		}, "/v1/agents/"},
		{"ICalSubscriptions.Get", func(c *chronary.Client, id string) error { return ignore(c.ICalSubscriptions.Get(ctx, id)) }, "/v1/ical-subscriptions/"},
		{"ICalSubscriptions.List", func(c *chronary.Client, id string) error {
			return ignore(c.ICalSubscriptions.List(ctx, id, nil).GetPage(ctx, 0))
		}, "/v1/agents/"},
		{"ICalSubscriptions.Update", func(c *chronary.Client, id string) error {
			return ignore(c.ICalSubscriptions.Update(ctx, id, &chronary.UpdateICalSubscriptionParams{}))
		}, "/v1/ical-subscriptions/"},
		{"ICalSubscriptions.Delete", func(c *chronary.Client, id string) error { return c.ICalSubscriptions.Delete(ctx, id) }, "/v1/ical-subscriptions/"},
		{"ICalSubscriptions.Sync", func(c *chronary.Client, id string) error { return c.ICalSubscriptions.Sync(ctx, id) }, "/v1/ical-subscriptions/"},
		{"ICalSubscriptions.Pause", func(c *chronary.Client, id string) error { return ignore(c.ICalSubscriptions.Pause(ctx, id)) }, "/v1/ical-subscriptions/"},
		{"ICalSubscriptions.Resume", func(c *chronary.Client, id string) error { return ignore(c.ICalSubscriptions.Resume(ctx, id)) }, "/v1/ical-subscriptions/"},
		{"Keys.Delete", func(c *chronary.Client, id string) error { return c.Keys.Delete(ctx, id) }, "/v1/keys/"},
		{"Scheduling.Get", func(c *chronary.Client, id string) error { return ignore(c.Scheduling.Get(ctx, id)) }, "/v1/scheduling/proposals/"},
		{"Scheduling.Respond", func(c *chronary.Client, id string) error {
			return ignore(c.Scheduling.Respond(ctx, id, &chronary.RespondToProposalParams{}))
		}, "/v1/scheduling/proposals/"},
		{"Scheduling.Resolve", func(c *chronary.Client, id string) error { return ignore(c.Scheduling.Resolve(ctx, id)) }, "/v1/scheduling/proposals/"},
		{"Scheduling.Cancel", func(c *chronary.Client, id string) error { return ignore(c.Scheduling.Cancel(ctx, id)) }, "/v1/scheduling/proposals/"},
		{"Webhooks.Get", func(c *chronary.Client, id string) error { return ignore(c.Webhooks.Get(ctx, id)) }, "/v1/webhooks/"},
		{"Webhooks.Update", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.Update(ctx, id, &chronary.UpdateWebhookParams{}))
		}, "/v1/webhooks/"},
		{"Webhooks.Delete", func(c *chronary.Client, id string) error { return c.Webhooks.Delete(ctx, id) }, "/v1/webhooks/"},
		{"Webhooks.ListDeliveries", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.ListDeliveries(ctx, id, nil))
		}, "/v1/webhooks/"},
		{"Webhooks.GetDelivery(webhookID)", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.GetDelivery(ctx, id, "whd_1"))
		}, "/v1/webhooks/"},
		{"Webhooks.GetDelivery(deliveryID)", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.GetDelivery(ctx, "whk_1", id))
		}, "/v1/webhooks/whk_1/deliveries/"},
		{"Webhooks.SendTest", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.SendTest(ctx, id, &chronary.SendTestWebhookParams{}))
		}, "/v1/webhooks/"},
		{"Webhooks.ReplayDelivery(webhookID)", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.ReplayDelivery(ctx, id, "whd_1"))
		}, "/v1/webhooks/"},
		{"Webhooks.ReplayDelivery(deliveryID)", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.ReplayDelivery(ctx, "whk_1", id))
		}, "/v1/webhooks/whk_1/deliveries/"},
		{"Webhooks.RotateSecret", func(c *chronary.Client, id string) error {
			return ignore(c.Webhooks.RotateSecret(ctx, id, nil))
		}, "/v1/webhooks/"},
	}
}

// recorder returns a client whose server records each request's escaped path.
func recorder(t *testing.T) (*chronary.Client, func() []string) {
	t.Helper()
	var mu sync.Mutex
	var paths []string
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.EscapedPath())
		mu.Unlock()
		testutil.RespondJSON(w, 200, map[string]interface{}{})
	}))
	return client, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), paths...)
	}
}

func TestPathParamsStayInOneSegment(t *testing.T) {
	for _, pc := range pathCalls() {
		t.Run(pc.name, func(t *testing.T) {
			client, sent := recorder(t)
			_ = pc.call(client, escapeID) // the stub answers {}; only the path matters
			got := sent()
			if len(got) != 1 {
				t.Fatalf("expected 1 request, got %v", got)
			}
			if !strings.HasPrefix(got[0], pc.prefix+escapeEncoded) {
				t.Errorf("path %q does not start with %q", got[0], pc.prefix+escapeEncoded)
			}
		})
	}
}

func TestPathParamsRefuseDotSegments(t *testing.T) {
	for _, pc := range pathCalls() {
		t.Run(pc.name, func(t *testing.T) {
			client, sent := recorder(t)
			for _, bad := range []string{"", ".", ".."} {
				err := pc.call(client, bad)
				if !errors.Is(err, chronary.ErrValidation) {
					t.Errorf("%q: expected a validation error, got %v", bad, err)
				}
				var apiErr *chronary.Error
				if errors.As(err, &apiErr) && apiErr.StatusCode != 0 {
					t.Errorf("%q: expected StatusCode 0 (nothing sent), got %d", bad, apiErr.StatusCode)
				}
			}
			if got := sent(); len(got) != 0 {
				t.Errorf("expected no request, got %v", got)
			}
		})
	}
}

func TestPathParamsIssueProofCases(t *testing.T) {
	client, sent := recorder(t)
	ctx := context.Background()
	_ = client.Events.DeleteByID(ctx, "../agents/agt_1")
	_, _ = client.Events.GetByID(ctx, "%2e%2e/keys")
	_ = client.Events.Delete(ctx, "../webhooks/whk_1#", "evt_x")
	_, _ = client.Events.GetByID(ctx, "evt_1?expand=all")
	want := []string{
		"/v1/events/..%2Fagents%2Fagt_1",
		"/v1/events/%252e%252e%2Fkeys",
		"/v1/calendars/..%2Fwebhooks%2Fwhk_1%23/events/evt_x",
		"/v1/events/evt_1%3Fexpand=all",
	}
	got := sent()
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("sent paths:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
