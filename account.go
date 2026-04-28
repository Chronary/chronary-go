package chronary

import (
	"context"
	"encoding/json"
	"net/http"
)

// AccountService handles GDPR portability + erasure operations.
type AccountService struct {
	service
}

// DataExport is the response shape for GET /v1/auth/export — every row this
// org owns, with encrypted fields decrypted in place.
//
// Sensitive fields (key hashes, password hashes, OTP hashes, claim revocation
// tokens, internal scheduling state) are excluded server-side.
type DataExport struct {
	ExportedAt    string                 `json:"exported_at"`
	FormatVersion string                 `json:"format_version"`
	Org           DataExportOrg          `json:"org"`
	Agents        []json.RawMessage      `json:"agents"`
	Calendars     []json.RawMessage      `json:"calendars"`
	Events        []json.RawMessage      `json:"events"`
	AvailabilityRules     []json.RawMessage `json:"availability_rules"`
	ICalSubscriptions     []json.RawMessage `json:"ical_subscriptions"`
	WebhookSubscriptions  []json.RawMessage `json:"webhook_subscriptions"`
	APIKeys              []json.RawMessage `json:"api_keys"`
	SchedulingProposals  []json.RawMessage `json:"scheduling_proposals"`
	ProposalSlots        []json.RawMessage `json:"proposal_slots"`
	ProposalResponses    []json.RawMessage `json:"proposal_responses"`
	UsageRecords         []json.RawMessage `json:"usage_records"`
	QuotaCounters        []json.RawMessage `json:"quota_counters"`
	TosAcceptances       []json.RawMessage `json:"tos_acceptances"`
	AccountClaimsInitiated []json.RawMessage `json:"account_claims_initiated"`
}

// DataExportOrg is the org block at the top of the export payload.
type DataExportOrg struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	Plan                  string  `json:"plan"`
	SignupSource          string  `json:"signup_source"`
	Status                string  `json:"status"`
	OAuthProvider         *string `json:"oauth_provider"`
	OAuthProviderID       *string `json:"oauth_provider_id"`
	EmailVerified         bool    `json:"email_verified"`
	OnboardingCompletedAt *string `json:"onboarding_completed_at"`
	AcceptedTermsVersion  *string `json:"accepted_terms_version"`
	AcceptedTermsAt       *string `json:"accepted_terms_at"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
}

// Export returns a JSON dump of every row this org owns (GDPR Art. 15 + 20
// portability / CCPA right-to-know / EU Data Act interoperability).
//
// AUTHENTICATION: This endpoint is JWT-only. It returns decrypted webhook
// secrets and iCal URLs that aren't normally accessible via API-key endpoints.
// Configure the client with a console JWT (cookie value or Bearer token from
// the console session) via WithAPIKey. API keys (chr_sk_* / chr_ak_*) will
// return 401.
//
// Most users should download via the console UI at console.chronary.ai/settings.
// This SDK method exists for programmatic compliance tooling holding a
// delegated JWT.
//
// Rate limit: 10 exports/hour/org.
func (s *AccountService) Export(ctx context.Context, opts ...RequestOption) (*DataExport, error) {
	var resp DataExport
	err := s.client.do(ctx, http.MethodGet, "/v1/auth/export", nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
