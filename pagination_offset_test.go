package chronary

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"
)

// bughunt-BH-094 (#2720): every list iterator hard-coded offset 0, so
// Collect / All ignored params.Offset and started at row 0, and
// BookingPages.List ignored params.Limit (always 50).

// offsetRecorder answers every list with one item and a total that ends the
// iteration after that page, and records the query of each request.
type offsetRecorder struct {
	mu      sync.Mutex
	queries []url.Values
}

func (o *offsetRecorder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	o.mu.Lock()
	o.queries = append(o.queries, q)
	o.mu.Unlock()
	offset, _ := strconv.Atoi(q.Get("offset"))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"data":[{}],"total":%d,"limit":%s,"offset":%d,"next_cursor":null}`, offset+1, q.Get("limit"), offset)
}

func (o *offsetRecorder) first(t *testing.T) url.Values {
	t.Helper()
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.queries) == 0 {
		t.Fatal("no request was sent")
	}
	return o.queries[0]
}

func TestListIteratorsStartAtParamsOffset(t *testing.T) {
	calID, agentID := "cal_1", "agt_1"
	cases := []struct {
		name    string
		collect func(ctx context.Context, c *Client) error
	}{
		{"Agents.List", func(ctx context.Context, c *Client) error {
			_, err := c.Agents.List(ctx, &ListAgentsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Agents.ListCalendars", func(ctx context.Context, c *Client) error {
			_, err := c.Agents.ListCalendars(ctx, "agt_1", &ListCalendarsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Agents.ListEvents", func(ctx context.Context, c *Client) error {
			_, err := c.Agents.ListEvents(ctx, "agt_1", &ListEventsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Agents.ListICalSubscriptions", func(ctx context.Context, c *Client) error {
			_, err := c.Agents.ListICalSubscriptions(ctx, "agt_1", &ListICalSubscriptionsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"ICalSubscriptions.List", func(ctx context.Context, c *Client) error {
			_, err := c.ICalSubscriptions.List(ctx, "agt_1", &ListICalSubscriptionsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"BookingPages.List", func(ctx context.Context, c *Client) error {
			_, err := c.BookingPages.List(ctx, &ListBookingPagesParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Calendars.List", func(ctx context.Context, c *Client) error {
			_, err := c.Calendars.List(ctx, &ListCalendarsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Events.List by calendar", func(ctx context.Context, c *Client) error {
			_, err := c.Events.List(ctx, &ListEventsParams{CalendarID: &calID, Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Events.List by agent", func(ctx context.Context, c *Client) error {
			_, err := c.Events.List(ctx, &ListEventsParams{AgentID: &agentID, Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Events.ListAll", func(ctx context.Context, c *Client) error {
			_, err := c.Events.ListAll(ctx, &ListOrgEventsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Scheduling.List", func(ctx context.Context, c *Client) error {
			_, err := c.Scheduling.List(ctx, &ListProposalsParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
		{"Webhooks.List", func(ctx context.Context, c *Client) error {
			_, err := c.Webhooks.List(ctx, &ListWebhooksParams{Offset: 30, Limit: 7}).Collect(ctx)
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := &offsetRecorder{}
			srv := httptest.NewServer(rec)
			defer srv.Close()
			client, err := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(0))
			if err != nil {
				t.Fatal(err)
			}
			if err := tc.collect(context.Background(), client); err != nil {
				t.Fatalf("Collect: %v", err)
			}
			q := rec.first(t)
			if got := q.Get("offset"); got != "30" {
				t.Errorf("first request offset = %q, want 30", got)
			}
			if got := q.Get("limit"); got != "7" {
				t.Errorf("first request limit = %q, want 7", got)
			}
		})
	}
}

// GetPage's explicit offset still wins over params.Offset.
func TestGetPageOffsetArgumentWins(t *testing.T) {
	rec := &offsetRecorder{}
	srv := httptest.NewServer(rec)
	defer srv.Close()
	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(0))

	if _, err := client.Webhooks.List(context.Background(), &ListWebhooksParams{Offset: 30}).GetPage(context.Background(), 5); err != nil {
		t.Fatal(err)
	}
	if got := rec.first(t).Get("offset"); got != "5" {
		t.Errorf("offset = %q, want 5", got)
	}
}

func TestListIteratorWithoutOffsetStartsAtZero(t *testing.T) {
	rec := &offsetRecorder{}
	srv := httptest.NewServer(rec)
	defer srv.Close()
	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(0))

	if _, err := client.BookingPages.List(context.Background(), nil).Collect(context.Background()); err != nil {
		t.Fatal(err)
	}
	q := rec.first(t)
	if got := q.Get("offset"); got != "0" {
		t.Errorf("offset = %q, want 0", got)
	}
	if got := q.Get("limit"); got != "50" {
		t.Errorf("limit = %q, want 50 (default)", got)
	}
}
