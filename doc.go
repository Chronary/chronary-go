// Package chronary provides a Go client for the Chronary calendar-as-a-service API.
//
// Create a client with your API key:
//
//	client, err := chronary.NewClient(
//	    chronary.WithAPIKey("chr_sk_live_..."),
//	)
//
// Or set the CHRONARY_API_KEY environment variable and call NewClient with no options:
//
//	client, err := chronary.NewClient()
//
// Access API resources through the client's service fields:
//
//	// Create a calendar
//	cal, err := client.Calendars.Create(ctx, &chronary.CreateCalendarParams{
//	    Name:     "Team Meetings",
//	    Timezone: "America/New_York",
//	})
//
//	// List events with pagination
//	for event, err := range client.Events.List(ctx, &chronary.ListEventsParams{
//	    CalendarID: chronary.String(cal.ID),
//	}).All(ctx) {
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Println(event.Title)
//	}
//
//	// Check availability
//	avail, err := client.Availability.Check(ctx, &chronary.CrossAgentAvailabilityParams{
//	    Agents: []string{"agt_abc", "agt_def"},
//	    Start:  "2026-04-15T09:00:00Z",
//	    End:    "2026-04-15T17:00:00Z",
//	})
package chronary
