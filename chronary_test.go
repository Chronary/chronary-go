package chronary

import (
	"os"
	"testing"
)

func TestNewClientRequiresAPIKey(t *testing.T) {
	os.Unsetenv("CHRONARY_API_KEY")
	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error when no API key is provided")
	}
}

func TestNewClientFromEnv(t *testing.T) {
	os.Setenv("CHRONARY_API_KEY", "chr_sk_env")
	defer os.Unsetenv("CHRONARY_API_KEY")

	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	if c.apiKey != "chr_sk_env" {
		t.Errorf("expected env key, got %s", c.apiKey)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	c, err := NewClient(
		WithAPIKey("test_key"),
		WithBaseURL("https://custom.api.com"),
		WithMaxRetries(5),
	)
	if err != nil {
		t.Fatal(err)
	}
	if c.apiKey != "test_key" {
		t.Errorf("expected test_key, got %s", c.apiKey)
	}
	if c.baseURL != "https://custom.api.com" {
		t.Errorf("expected custom URL, got %s", c.baseURL)
	}
	if c.maxRetries != 5 {
		t.Errorf("expected 5 retries, got %d", c.maxRetries)
	}
}

func TestNewClientTrimsTrailingSlash(t *testing.T) {
	c, err := NewClient(
		WithAPIKey("test_key"),
		WithBaseURL("https://api.chronary.ai/"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != "https://api.chronary.ai" {
		t.Errorf("expected trimmed URL, got %s", c.baseURL)
	}
}

func TestNewClientDefaults(t *testing.T) {
	c, err := NewClient(WithAPIKey("test_key"))
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != DefaultBaseURL {
		t.Errorf("expected default base URL, got %s", c.baseURL)
	}
	if c.maxRetries != DefaultMaxRetries {
		t.Errorf("expected default retries, got %d", c.maxRetries)
	}
	if c.timeout != DefaultTimeout {
		t.Errorf("expected default timeout, got %s", c.timeout)
	}
}

func TestNewClientServicesInitialized(t *testing.T) {
	c, err := NewClient(WithAPIKey("test_key"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Agents == nil {
		t.Error("Agents service is nil")
	}
	if c.Calendars == nil {
		t.Error("Calendars service is nil")
	}
	if c.Events == nil {
		t.Error("Events service is nil")
	}
	if c.Availability == nil {
		t.Error("Availability service is nil")
	}
	if c.Webhooks == nil {
		t.Error("Webhooks service is nil")
	}
	if c.ICalSubscriptions == nil {
		t.Error("ICalSubscriptions service is nil")
	}
	if c.Usage == nil {
		t.Error("Usage service is nil")
	}
}
