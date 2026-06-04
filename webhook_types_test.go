package chronary

import "testing"

// The full webhook event-type set must stay in lockstep with the server's
// WEBHOOK_EVENT_TYPES and the TypeScript / Python SDK unions (18 events).
func TestWebhookEventTypesParity(t *testing.T) {
	want := []WebhookEventType{
		"agent.created", "agent.updated",
		"event.created", "event.updated", "event.deleted", "event.started", "event.ended", "event.reminder",
		"event.hold_created", "event.hold_expired", "event.hold_released", "event.hold_confirmed",
		"proposal.created", "proposal.responded", "proposal.confirmed", "proposal.expired", "proposal.cancelled",
		"webhook.deactivated",
	}

	if len(WebhookEventTypes) != len(want) {
		t.Fatalf("expected %d webhook event types, got %d", len(want), len(WebhookEventTypes))
	}

	got := make(map[WebhookEventType]bool, len(WebhookEventTypes))
	for _, e := range WebhookEventTypes {
		got[e] = true
	}
	for _, w := range want {
		if !got[w] {
			t.Errorf("WebhookEventTypes is missing %q", w)
		}
	}

	// webhook.deactivated was the drift that prompted this — guard it explicitly.
	if !got[WebhookEventWebhookDeactivated] {
		t.Error("WebhookEventTypes must include webhook.deactivated")
	}
	if WebhookEventWebhookDeactivated != "webhook.deactivated" {
		t.Errorf("WebhookEventWebhookDeactivated = %q, want webhook.deactivated", WebhookEventWebhookDeactivated)
	}
}
