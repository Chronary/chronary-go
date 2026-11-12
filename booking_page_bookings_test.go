package chronary_test

import (
	"context"
	"net/http"
	"testing"

	chronary "github.com/Chronary/chronary-go"
	"github.com/Chronary/chronary-go/internal/testutil"
)

// #1202: booking page calendar filter, bookings list and link rotation.

func bookingPageJSON(slug string) map[string]interface{} {
	return map[string]interface{}{
		"id": "bkp_1", "org_id": "org_1", "calendar_id": "cal_1", "agent_id": nil,
		"slug": slug, "title": "Intro", "description": nil,
		"duration_minutes": 30, "buffer_minutes": 0, "window_days": 14, "min_notice_minutes": 0,
		"timezone": "UTC", "availability_constraints": nil, "active": true, "bookings_count": 1,
		"booking_url": "https://api.chronary.ai/book/" + slug,
		"created_at":  "2026-07-14T10:00:00Z", "updated_at": "2026-07-14T10:00:00Z",
	}
}

func TestBookingPageList_CalendarID(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertPath(t, r, "/v1/booking-pages")
		if got := r.URL.Query().Get("calendar_id"); got != "cal_1" {
			t.Errorf("calendar_id = %q", got)
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []interface{}{bookingPageJSON("s1")}, "total": 1, "limit": 50, "offset": 0,
		})
	}))

	page, err := client.BookingPages.List(context.Background(), &chronary.ListBookingPagesParams{
		CalendarID: chronary.String("cal_1"),
	}).GetPage(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Data) != 1 || page.Data[0].ID != "bkp_1" {
		t.Errorf("unexpected page %+v", page.Data)
	}
}

func TestBookingPageListBookings(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/booking-pages/bkp_1/bookings")
		q := r.URL.Query()
		if q.Get("limit") != "10" || q.Get("cursor") != "abc" {
			t.Errorf("unexpected query %s", r.URL.RawQuery)
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []map[string]interface{}{{
				"id": "evt_1", "booking_page_id": "bkp_1", "calendar_id": "cal_1",
				"start_time": "2026-07-15T16:00:00Z", "end_time": "2026-07-15T16:30:00Z",
				"status": "confirmed", "invitee_name": "Ada", "invitee_email": "ada@example.com",
				"invitee_notes": nil, "created_at": "2026-07-14T10:00:00Z",
			}},
			"next_cursor": "next",
		})
	}))

	list, err := client.BookingPages.ListBookings(context.Background(), "bkp_1", &chronary.ListBookingsParams{
		Limit:  10,
		Cursor: chronary.String("abc"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != "evt_1" || list.Data[0].InviteeEmail == nil || *list.Data[0].InviteeEmail != "ada@example.com" {
		t.Errorf("unexpected bookings %+v", list.Data)
	}
	if list.Data[0].InviteeNotes != nil {
		t.Errorf("invitee_notes = %v, want nil", *list.Data[0].InviteeNotes)
	}
	if list.NextCursor == nil || *list.NextCursor != "next" {
		t.Errorf("next_cursor = %v", list.NextCursor)
	}
}

func TestBookingPageListBookings_NilParams(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertPath(t, r, "/v1/booking-pages/bkp_1/bookings")
		if r.URL.RawQuery != "" {
			t.Errorf("unexpected query %s", r.URL.RawQuery)
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{"data": []interface{}{}, "next_cursor": nil})
	}))

	list, err := client.BookingPages.ListBookings(context.Background(), "bkp_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 0 || list.NextCursor != nil {
		t.Errorf("unexpected list %+v", list)
	}
}

func TestBookingPageRotateSlug(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/booking-pages/bkp_1/rotate-slug")
		testutil.RespondJSON(w, 200, bookingPageJSON("newslug"))
	}))

	page, err := client.BookingPages.RotateSlug(context.Background(), "bkp_1")
	if err != nil {
		t.Fatal(err)
	}
	if page.Slug != "newslug" || page.BookingURL != "https://api.chronary.ai/book/newslug" {
		t.Errorf("unexpected page %+v", page)
	}
}

func TestBookingPageErrors(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.RespondJSON(w, 404, map[string]interface{}{
			"error": map[string]interface{}{"type": "not_found", "message": "Booking page not found"},
		})
	}))

	if _, err := client.BookingPages.ListBookings(context.Background(), "bkp_x", nil); err == nil {
		t.Error("ListBookings: expected an error")
	}
	if _, err := client.BookingPages.RotateSlug(context.Background(), "bkp_x"); err == nil {
		t.Error("RotateSlug: expected an error")
	}
}
