package chronary

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

// bughunt-BH-080 (#2691): with NOC_QUERY_CAPS on, every list rejects
// offset > 10000 with 400 validation_error. The iterators paged by offset
// only, so Collect / All failed partway past 10,000 rows, even on lists that
// return next_cursor. They now follow next_cursor where the list has one.

const (
	cappedMaxOffset = 10000
	cappedRows      = 12000
)

type cappedAPI struct {
	withCursor bool
	mu         sync.Mutex
	offsets    []int
	cursors    []string
}

func encIndex(i int) string { return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(i))) }

func (a *cappedAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	cursor, hasCursor := q.Get("cursor"), q.Has("cursor")
	writeErr := func(msg string) {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error":{"type":"validation_error","message":%q}}`, msg)
	}
	if offset > cappedMaxOffset {
		writeErr("offset must be at most 10000")
		return
	}
	if hasCursor && (!a.withCursor || offset > 0) {
		writeErr("cursor and offset cannot be combined")
		return
	}
	a.mu.Lock()
	a.offsets = append(a.offsets, offset)
	if hasCursor {
		a.cursors = append(a.cursors, cursor)
	}
	a.mu.Unlock()

	start := offset
	if hasCursor {
		raw, _ := base64.StdEncoding.DecodeString(cursor)
		start, _ = strconv.Atoi(string(raw))
	}
	stop := start + limit
	if stop > cappedRows {
		stop = cappedRows
	}
	data := make([]map[string]any, 0, stop-start)
	for i := start; i < stop; i++ {
		data = append(data, map[string]any{"id": fmt.Sprintf("row_%d", i)})
	}
	body := map[string]any{"data": data, "total": cappedRows, "limit": limit, "offset": offset}
	if a.withCursor {
		if stop < cappedRows {
			body["next_cursor"] = encIndex(stop)
		} else {
			body["next_cursor"] = nil
		}
	}
	_ = json.NewEncoder(w).Encode(body)
}

func newCappedClient(t *testing.T, api *cappedAPI) *Client {
	t.Helper()
	srv := httptest.NewServer(api)
	t.Cleanup(srv.Close)
	client, err := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestCursorPagingEventsListAllPastOffsetCap(t *testing.T) {
	api := &cappedAPI{withCursor: true}
	client := newCappedClient(t, api)

	items, err := client.Events.ListAll(context.Background(), &ListOrgEventsParams{Limit: 100}).Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	if len(items) != cappedRows {
		t.Fatalf("expected %d events, got %d", cappedRows, len(items))
	}
	seen := map[string]bool{}
	for _, e := range items {
		seen[e.ID] = true
	}
	if len(seen) != cappedRows {
		t.Errorf("expected %d distinct events, got %d", cappedRows, len(seen))
	}
	for _, off := range api.offsets {
		if off != 0 {
			t.Fatalf("expected every request at offset 0 (cursor paging), saw offset %d", off)
		}
	}
	if len(api.cursors) != cappedRows/100-1 {
		t.Errorf("expected %d cursor requests, got %d", cappedRows/100-1, len(api.cursors))
	}
}

func TestCursorPagingAgentsAndCalendarsPastOffsetCap(t *testing.T) {
	client := newCappedClient(t, &cappedAPI{withCursor: true})
	ctx := context.Background()

	agents, err := client.Agents.List(ctx, &ListAgentsParams{Limit: 200}).Collect(ctx)
	if err != nil || len(agents) != cappedRows {
		t.Fatalf("Agents.List: %d rows, err %v", len(agents), err)
	}
	cals, err := client.Calendars.List(ctx, &ListCalendarsParams{Limit: 200}).Collect(ctx)
	if err != nil || len(cals) != cappedRows {
		t.Fatalf("Calendars.List: %d rows, err %v", len(cals), err)
	}
	agentCals, err := client.Agents.ListCalendars(ctx, "agt_1", &ListCalendarsParams{Limit: 200}).Collect(ctx)
	if err != nil || len(agentCals) != cappedRows {
		t.Fatalf("Agents.ListCalendars: %d rows, err %v", len(agentCals), err)
	}
}

func TestCursorPagingStartsFromCallerCursor(t *testing.T) {
	client := newCappedClient(t, &cappedAPI{withCursor: true})
	ctx := context.Background()

	agents, err := client.Agents.List(ctx, &ListAgentsParams{Limit: 200, Cursor: String(encIndex(11000))}).Collect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(agents) != 1000 || agents[0].ID != "row_11000" {
		t.Fatalf("expected 1000 agents from row_11000, got %d starting %v", len(agents), agents)
	}
}

func TestCursorPagingGetPageExposesNextCursor(t *testing.T) {
	client := newCappedClient(t, &cappedAPI{withCursor: true})
	ctx := context.Background()
	iter := client.Events.ListAll(ctx, &ListOrgEventsParams{Limit: 100})

	first, err := iter.GetPage(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if first.NextCursor == nil || *first.NextCursor != encIndex(100) || !first.HasMore {
		t.Fatalf("first page: NextCursor %v, HasMore %v", first.NextCursor, first.HasMore)
	}
	second, err := iter.GetPageAfter(ctx, *first.NextCursor)
	if err != nil {
		t.Fatal(err)
	}
	if second.Data[0].ID != "row_100" {
		t.Errorf("expected row_100, got %s", second.Data[0].ID)
	}
	last, err := iter.GetPageAfter(ctx, encIndex(11900))
	if err != nil {
		t.Fatal(err)
	}
	if last.NextCursor != nil || last.HasMore {
		t.Errorf("last page: NextCursor %v, HasMore %v", last.NextCursor, last.HasMore)
	}
}

func TestOffsetOnlyListStillPagesByOffset(t *testing.T) {
	api := &cappedAPI{withCursor: false}
	client := newCappedClient(t, api)
	ctx := context.Background()

	items, err := client.Events.List(ctx, &ListEventsParams{CalendarID: String("cal_1"), Limit: 2000}).Collect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != cappedRows {
		t.Fatalf("expected %d events, got %d", cappedRows, len(items))
	}
	if len(api.cursors) != 0 {
		t.Errorf("an offset-only list must never be sent a cursor, got %v", api.cursors)
	}
	want := []int{0, 2000, 4000, 6000, 8000, 10000}
	if fmt.Sprint(api.offsets) != fmt.Sprint(want) {
		t.Errorf("offsets %v, want %v", api.offsets, want)
	}
}

// Review of #2748: a list that hands back the cursor just sent must not loop.
func TestCursorPagingStopsOnRepeatedCursor(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		fmt.Fprint(w, `{"data":[{"id":"a"}],"total":5,"limit":1,"offset":0,"next_cursor":"same"}`)
	}))
	defer srv.Close()
	client, _ := NewClient(WithAPIKey("test_key"), WithBaseURL(srv.URL), WithMaxRetries(0))

	_, err := client.Agents.List(context.Background(), &ListAgentsParams{Limit: 1}).Collect(context.Background())
	if err == nil {
		t.Fatal("expected an error for a repeated cursor")
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("expected 2 requests, got %d", got)
	}
}
