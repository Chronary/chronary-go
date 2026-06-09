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

func TestAgentCreate(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "POST")
		testutil.AssertPath(t, r, "/v1/agents")
		testutil.AssertHasAuth(t, r)
		testutil.AssertHasIdempotencyKey(t, r)

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "Bot" {
			t.Errorf("expected name Bot, got %v", body["name"])
		}

		testutil.RespondJSON(w, 201, map[string]interface{}{
			"id": "agt_1", "name": "Bot", "type": "ai", "status": "active",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
		})
	}))

	agent, err := client.Agents.Create(context.Background(), &chronary.CreateAgentParams{
		Name: "Bot",
		Type: chronary.AgentTypeAI,
	})
	if err != nil {
		t.Fatal(err)
	}
	if agent.ID != "agt_1" {
		t.Errorf("expected agt_1, got %s", agent.ID)
	}
	if agent.Name != "Bot" {
		t.Errorf("expected Bot, got %s", agent.Name)
	}
}

func TestAgentGet(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.AssertPath(t, r, "/v1/agents/agt_1")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "agt_1", "orgId": "org_1", "name": "Bot", "type": "ai", "status": "active",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
		})
	}))

	agent, err := client.Agents.Get(context.Background(), "agt_1")
	if err != nil {
		t.Fatal(err)
	}
	if agent.ID != "agt_1" {
		t.Errorf("expected agt_1, got %s", agent.ID)
	}
	// Guard the live response shape (camelCase): these deserialized as zero
	// values when the struct tags were snake_case.
	if agent.OrgID != "org_1" {
		t.Errorf("expected orgId org_1, got %q", agent.OrgID)
	}
	if agent.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to deserialize from camelCase createdAt, got zero time")
	}
	if agent.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to deserialize from camelCase updatedAt, got zero time")
	}
}

func TestAgentUpdate(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "PATCH")
		testutil.AssertPath(t, r, "/v1/agents/agt_1")
		testutil.AssertHasIdempotencyKey(t, r)
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"id": "agt_1", "name": "Updated Bot", "type": "ai", "status": "active",
			"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z",
		})
	}))

	agent, err := client.Agents.Update(context.Background(), "agt_1", &chronary.UpdateAgentParams{
		Name: chronary.String("Updated Bot"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if agent.Name != "Updated Bot" {
		t.Errorf("expected Updated Bot, got %s", agent.Name)
	}
}

func TestAgentDelete(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "DELETE")
		testutil.AssertPath(t, r, "/v1/agents/agt_1")
		testutil.AssertHasIdempotencyKey(t, r)
		w.WriteHeader(204)
	}))

	err := client.Agents.Delete(context.Background(), "agt_1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAgentList(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertMethod(t, r, "GET")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "agt_1", "name": "Bot 1", "type": "ai", "status": "active",
					"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z"},
				{"id": "agt_2", "name": "Bot 2", "type": "human", "status": "active",
					"createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z"},
			},
			"total": 2,
		})
	}))

	items, err := client.Agents.List(context.Background(), nil).Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 agents, got %d", len(items))
	}
}

func TestAgentGetNotFound(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.RespondError(w, 404, "not_found", "Agent not found")
	}))

	_, err := client.Agents.Get(context.Background(), "agt_missing")
	if !errors.Is(err, chronary.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAgentListCalendars(t *testing.T) {
	client, _ := testutil.Setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testutil.AssertPath(t, r, "/v1/agents/agt_1/calendars")
		testutil.RespondJSON(w, 200, map[string]interface{}{
			"data":  []map[string]interface{}{{"id": "cal_1", "name": "Work", "timezone": "UTC", "createdAt": "2026-04-14T00:00:00Z", "updatedAt": "2026-04-14T00:00:00Z"}},
			"total": 1,
		})
	}))

	items, err := client.Agents.ListCalendars(context.Background(), "agt_1", nil).Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 calendar, got %d", len(items))
	}
}
