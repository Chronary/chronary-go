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
