package chronary

import "time"

// String returns a pointer to s.
func String(s string) *string { return &s }

// Int returns a pointer to i.
func Int(i int) *int { return &i }

// Int64 returns a pointer to i.
func Int64(i int64) *int64 { return &i }

// Float64 returns a pointer to f.
func Float64(f float64) *float64 { return &f }

// Bool returns a pointer to b.
func Bool(b bool) *bool { return &b }

// Time returns a pointer to t.
func Time(t time.Time) *time.Time { return &t }

// StringValue returns the value of a string pointer, or empty string if nil.
func StringValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// IntValue returns the value of an int pointer, or 0 if nil.
func IntValue(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// BoolValue returns the value of a bool pointer, or false if nil.
func BoolValue(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

// Reminders returns a reminder schedule for the DefaultReminders / Reminders
// fields of the calendar and event create/update params: offsets in minutes
// before an event's start_time, e.g. Reminders(10, 1440). The result is never
// a nil slice, so Reminders() with no arguments sends [] (no reminders), the
// same as NoReminders().
func Reminders(minutes ...int) *[]int {
	r := make([]int, len(minutes))
	copy(r, minutes)
	return &r
}

// NoReminders returns an explicitly empty reminder schedule. It sends [], which
// means "no reminders" (on a calendar: no default reminders for its events).
// A nil field would be omitted and the event would inherit the default instead.
func NoReminders() *[]int {
	r := []int{}
	return &r
}

// InheritReminders returns a pointer to a nil slice, which sends an explicit
// JSON null. On an update it reverts an event to inheriting its calendar's
// default reminders, or a calendar to the system default ([10]). On a create it
// is the same as leaving the field nil.
func InheritReminders() *[]int {
	return new([]int)
}
