package chronary

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default Chronary API base URL.
	DefaultBaseURL = "https://api.chronary.ai"
	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 30 * time.Second
	// DefaultMaxRetries is the default number of retries for failed requests.
	DefaultMaxRetries = 2
)

// Client is the Chronary API client. Access resources via the service fields.
type Client struct {
	Agents            *AgentService
	AgentAuth         *AgentAuthService
	Calendars         *CalendarService
	Events            *EventService
	Availability      *AvailabilityService
	Webhooks          *WebhookService
	ICalSubscriptions *ICalSubscriptionService
	Usage             *UsageService
	Scheduling        *SchedulingService
	Keys              *KeysService
	Feedback          *FeedbackService
	Plans             *PlansService
	Account           *AccountService

	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	maxRetries int
}

// NewClient creates a new Chronary API client.
//
// If no API key is provided via WithAPIKey, it reads from CHRONARY_API_KEY env var.
// To construct a client for unauthenticated endpoints (AgentAuth.SignUp, Plans.List),
// pass WithAnonymous() — the API key check is then skipped.
func NewClient(opts ...ClientOption) (*Client, error) {
	cfg := clientConfig{
		baseURL:    DefaultBaseURL,
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.apiKey == "" && !cfg.anonymous {
		cfg.apiKey = os.Getenv("CHRONARY_API_KEY")
	}
	if cfg.apiKey == "" && !cfg.anonymous {
		return nil, fmt.Errorf("chronary: API key is required (set CHRONARY_API_KEY, use WithAPIKey, or WithAnonymous for sign-up/plans)")
	}

	c := &Client{
		apiKey:     cfg.apiKey,
		baseURL:    strings.TrimRight(cfg.baseURL, "/"),
		httpClient: cfg.httpClient,
		timeout:    cfg.timeout,
		maxRetries: cfg.maxRetries,
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: c.timeout}
	}

	s := service{client: c}
	c.Agents = &AgentService{s}
	c.AgentAuth = &AgentAuthService{s}
	c.Calendars = &CalendarService{s}
	c.Events = &EventService{s}
	c.Availability = &AvailabilityService{s}
	c.Webhooks = &WebhookService{s}
	c.ICalSubscriptions = &ICalSubscriptionService{s}
	c.Usage = &UsageService{s}
	c.Scheduling = &SchedulingService{s}
	c.Keys = &KeysService{s}
	c.Feedback = &FeedbackService{s}
	c.Plans = &PlansService{s}
	c.Account = &AccountService{s}

	return c, nil
}
