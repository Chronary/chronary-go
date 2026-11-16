package chronary

import "testing"

func TestString(t *testing.T) {
	s := String("hello")
	if *s != "hello" {
		t.Errorf("expected hello, got %s", *s)
	}
}

func TestStringValue(t *testing.T) {
	if got := StringValue(nil); got != "" {
		t.Errorf("expected empty, got %s", got)
	}
	s := "hello"
	if got := StringValue(&s); got != "hello" {
		t.Errorf("expected hello, got %s", got)
	}
}

func TestInt(t *testing.T) {
	i := Int(42)
	if *i != 42 {
		t.Errorf("expected 42, got %d", *i)
	}
}

func TestIntValue(t *testing.T) {
	if got := IntValue(nil); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
	i := 42
	if got := IntValue(&i); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestBool(t *testing.T) {
	b := Bool(true)
	if *b != true {
		t.Error("expected true")
	}
}

func TestBoolValue(t *testing.T) {
	if got := BoolValue(nil); got != false {
		t.Error("expected false")
	}
	b := true
	if got := BoolValue(&b); got != true {
		t.Error("expected true")
	}
}

func TestSetTimezone(t *testing.T) {
	tz := SetTimezone("America/Los_Angeles")
	if tz == nil || *tz == nil || **tz != "America/Los_Angeles" {
		t.Errorf("expected pointer to America/Los_Angeles, got %v", tz)
	}
}

func TestClearTimezone(t *testing.T) {
	tz := ClearTimezone()
	if tz == nil {
		t.Fatal("expected a non-nil outer pointer")
	}
	if *tz != nil {
		t.Errorf("expected a nil inner pointer, got %v", **tz)
	}
}
