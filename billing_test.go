package chronary_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	chronary "github.com/Chronary/chronary-go"
	"github.com/Chronary/chronary-go/internal/testutil"
)

func upgradeRequestJSON(status string) map[string]interface{} {
	return map[string]interface{}{
		"id": "upr_1", "status": status, "plan": "pro", "interval": "month",
		"reason": "Hit the event limit", "feature": nil, "metric": "events",
		"requested_via": "api", "requested_by_key_prefix": "chr_sk_abcdefghijkl",
		"approval_url": "https://console.test/settings/billing/requests/upr_1",
		"created_at":   "2026-10-19T00:00:00Z", "expires_at": "2026-10-26T00:00:00Z",
		"completed_at": nil, "declined_at": nil, "canceled_at": nil,
	}
}

func TestBillingGet(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, http.MethodGet)
		testutil.AssertPath(t, r, "/v1/billing")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"plan": "free", "plan_source": nil, "plan_expires_at": nil, "status": nil,
			"interval": nil, "current_period_end": nil, "cancel_at": nil, "cancel_at_period_end": false,
			"renews_at": nil, "price": nil, "has_billing_account": false, "self_serve_upgrade": true,
			"upgrade_options": []interface{}{
				map[string]interface{}{"plan": "pro", "interval": "month", "amount": 2900, "currency": "usd", "savings_vs_monthly": 0},
			},
			"upgrade_url": "https://console.test/settings/billing", "manage_url": "https://console.test/settings/billing",
			"contact_url": nil, "pending_upgrade_request": upgradeRequestJSON("pending"), "offer": nil,
			"past_due": map[string]interface{}{
				"amount_due": 1581, "currency": "usd", "failed_at": "2026-10-23T00:00:00.000Z",
				"next_payment_attempt": nil, "hosted_invoice_url": "https://invoice.test/pay", "attempt_count": 4,
			},
			"credit": map[string]interface{}{"amount": 26100, "currency": "usd"},
		})
	}))

	b, err := client.Billing.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if b.Plan != "free" || !b.SelfServeUpgrade || len(b.UpgradeOptions) != 1 || b.UpgradeOptions[0].Amount != 2900 {
		t.Fatalf("unexpected summary: %#v", b)
	}
	if b.PendingUpgradeRequest == nil || b.PendingUpgradeRequest.ID != "upr_1" {
		t.Fatalf("expected pending request: %#v", b.PendingUpgradeRequest)
	}
	if b.PastDue == nil || b.PastDue.AmountDue != 1581 || b.PastDue.NextPaymentAttempt != nil || b.PastDue.AttemptCount != 4 {
		t.Fatalf("unexpected past_due: %#v", b.PastDue)
	}
	if b.Credit == nil || b.Credit.Amount != 26100 || b.Credit.Currency != "usd" {
		t.Fatalf("unexpected credit: %#v", b.Credit)
	}
	if b.Status != nil || b.UpgradeURL == nil {
		t.Fatalf("nullable handling wrong: %#v", b)
	}
}

func TestBillingUpgradeRequestLifecycle(t *testing.T) {
	step := 0
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		step++
		switch step {
		case 1:
			testutil.AssertMethod(t, r, http.MethodPost)
			testutil.AssertPath(t, r, "/v1/billing/upgrade-requests")
			if got := r.Header.Get("Idempotency-Key"); got != "idem-1" {
				t.Errorf("Idempotency-Key = %q, want idem-1", got)
			}
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["interval"] != "year" || body["reason"] != "need more" {
				t.Errorf("unexpected body: %#v", body)
			}
			resp := upgradeRequestJSON("pending")
			resp["owner_notification"] = "queued"
			resp["reused"] = false
			testutil.RespondJSON(w, 201, resp)
		case 2:
			testutil.AssertMethod(t, r, http.MethodGet)
			testutil.AssertPath(t, r, "/v1/billing/upgrade-requests/upr_1")
			testutil.RespondJSON(w, 200, upgradeRequestJSON("completed"))
		case 3:
			testutil.AssertMethod(t, r, http.MethodPost)
			testutil.AssertPath(t, r, "/v1/billing/upgrade-requests/upr_1/cancel")
			testutil.RespondJSON(w, 200, upgradeRequestJSON("canceled"))
		}
	}))

	ctx := context.Background()
	req, err := client.Billing.CreateUpgradeRequest(ctx, &chronary.CreateUpgradeRequestParams{Interval: "year", Reason: "need more"}, chronary.WithIdempotencyKey("idem-1"))
	if err != nil || req.ID != "upr_1" || req.ApprovalURL == nil || req.OwnerNotification != "queued" {
		t.Fatalf("create: %#v, %v", req, err)
	}
	got, err := client.Billing.GetUpgradeRequest(ctx, "upr_1")
	if err != nil || got.Status != "completed" {
		t.Fatalf("get: %#v, %v", got, err)
	}
	canceled, err := client.Billing.CancelUpgradeRequest(ctx, "upr_1")
	if err != nil || canceled.Status != "canceled" {
		t.Fatalf("cancel: %#v, %v", canceled, err)
	}
}

func TestBillingCreateUpgradeRequestNilParams(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body) != 0 {
			t.Errorf("expected empty body, got %#v", body)
		}
		resp := upgradeRequestJSON("pending")
		resp["reused"] = true
		testutil.RespondJSON(w, 200, resp)
	}))
	req, err := client.Billing.CreateUpgradeRequest(context.Background(), nil)
	if err != nil || !req.Reused {
		t.Fatalf("create: %#v, %v", req, err)
	}
}

func TestPlanRequiredErrorExposesUpgradeURL(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.RespondJSON(w, 403, map[string]interface{}{"error": map[string]interface{}{
			"type": "plan_required", "message": "Pro required", "request_id": "req_1",
			"required_plan": "pro", "feature": "scoped_keys",
			"upgrade_url": "https://console.test/settings/billing", "docs_url": "https://docs.test/plans",
			"contact_url": "https://chronary.ai/contact",
		}})
	}))

	_, err := client.Billing.Get(context.Background())
	if !errors.Is(err, chronary.ErrPlanRequired) {
		t.Fatalf("expected ErrPlanRequired, got %v", err)
	}
	var apiErr *chronary.Error
	if !errors.As(err, &apiErr) {
		t.Fatal("expected *chronary.Error")
	}
	if apiErr.StatusCode != 403 || apiErr.UpgradeURL != "https://console.test/settings/billing" ||
		apiErr.DocsURL != "https://docs.test/plans" || apiErr.RequiredPlan != "pro" ||
		apiErr.Feature != "scoped_keys" || apiErr.ContactURL != "https://chronary.ai/contact" {
		t.Fatalf("unexpected fields: %#v", apiErr)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(apiErr.Body, &raw); err != nil || raw["required_plan"] != "pro" {
		t.Fatalf("raw body not preserved: %s (%v)", apiErr.Body, err)
	}
}

func TestQuotaExceededKeepsStatusAndFields(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.RespondJSON(w, 429, map[string]interface{}{"error": map[string]interface{}{
			"type": "quota_exceeded", "message": "Quota hit", "request_id": "req_2",
			"metric": "api_calls", "limit": 50000, "upgrade_url": "https://console.test/settings/billing",
		}})
	}))

	_, err := client.Billing.Get(context.Background())
	if !errors.Is(err, chronary.ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
	var apiErr *chronary.Error
	if !errors.As(err, &apiErr) {
		t.Fatal("expected *chronary.Error")
	}
	if apiErr.StatusCode != 429 || apiErr.Metric != "api_calls" || apiErr.Limit != 50000 ||
		apiErr.UpgradeURL != "https://console.test/settings/billing" {
		t.Fatalf("unexpected fields: %#v", apiErr)
	}
}

func TestDeveloperPlanFields(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/plans" {
			testutil.RespondJSON(w, 200, map[string]interface{}{"plans": []map[string]interface{}{
				{"id": "developer", "name": "Developer", "tagline": "t", "price": 900, "annual_price": 9000,
					"currency": "usd", "limits": nil, "display_features": []string{}, "recommended": false},
				{"id": "some_future_plan", "name": "F", "tagline": "t", "price": nil, "currency": nil,
					"limits": nil, "display_features": []string{}, "recommended": false},
			}})
			return
		}
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"plan": "pro", "cancel_at_period_end": false, "has_billing_account": true, "self_serve_upgrade": false,
			"pending_plan": "developer", "pending_plan_interval": "month", "pending_plan_at": "2026-12-01T00:00:00Z",
			"upgrade_options": []map[string]interface{}{
				{"plan": "developer", "interval": "month", "amount": 900, "currency": "usd", "savings_vs_monthly": 0, "via": "switch_plan"},
			},
			"manage_url": "https://console.test/settings/billing",
		})
	}))
	plans, err := client.Plans.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if plans.Plans[0].ID != "developer" || plans.Plans[0].AnnualPrice == nil || *plans.Plans[0].AnnualPrice != 9000 {
		t.Fatalf("developer plan not decoded: %#v", plans.Plans[0])
	}
	if plans.Plans[1].ID != "some_future_plan" || plans.Plans[1].AnnualPrice != nil {
		t.Fatalf("unknown plan id must decode: %#v", plans.Plans[1])
	}
	s, err := client.Billing.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s.PendingPlan == nil || *s.PendingPlan != "developer" || s.PendingPlanInterval == nil ||
		*s.PendingPlanInterval != "month" || s.PendingPlanAt == nil {
		t.Fatalf("pending plan not decoded: %#v", s)
	}
	if len(s.UpgradeOptions) != 1 || s.UpgradeOptions[0].Plan != "developer" || s.UpgradeOptions[0].Via != "switch_plan" {
		t.Fatalf("upgrade option not decoded: %#v", s.UpgradeOptions)
	}
}

func TestQuotaErrorExposesUpgradePlan(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.RespondJSON(w, 429, map[string]interface{}{"error": map[string]interface{}{
			"type": "quota_exceeded", "message": "limit", "metric": "api_calls", "limit": 50000,
			"upgrade_url": "https://console.test/settings/billing", "upgrade_plan": "developer",
		}})
	}))
	_, err := client.Billing.Get(context.Background())
	var apiErr *chronary.Error
	if !errors.As(err, &apiErr) || apiErr.UpgradePlan != "developer" {
		t.Fatalf("expected UpgradePlan developer, got %#v (%v)", apiErr, err)
	}
}
