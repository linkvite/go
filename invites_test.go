package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestInvitesList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/invites" {
			t.Errorf("expected path /v1/invites, got %s", r.URL.Path)
		}

		data := []map[string]any{
			{
				"id":            "inv_123",
				"collection_id": "coll_456",
				"role":          "viewer",
				"collection": map[string]any{
					"id":   "coll_456",
					"name": "Shared Collection",
				},
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	invites, err := client.Invites.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(invites) != 1 {
		t.Errorf("expected 1 invite, got %d", len(invites))
	}
	if invites[0].ID != "inv_123" {
		t.Errorf("expected ID 'inv_123', got '%s'", invites[0].ID)
	}
}

func TestInvitesGet(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/invites/inv_123" {
			t.Errorf("expected path /v1/invites/inv_123, got %s", r.URL.Path)
		}

		data := map[string]any{
			"id":            "inv_123",
			"collection_id": "coll_456",
			"role":          "viewer",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	invite, err := client.Invites.Get(context.Background(), "inv_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invite.ID != "inv_123" {
		t.Errorf("expected ID 'inv_123', got '%s'", invite.ID)
	}
}

func TestInvitesAccept(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/invites/inv_123/accept" {
			t.Errorf("expected path /v1/invites/inv_123/accept, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Invites.Accept(context.Background(), "inv_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvitesDecline(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/invites/inv_123/decline" {
			t.Errorf("expected path /v1/invites/inv_123/decline, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Invites.Decline(context.Background(), "inv_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvitesBatchAccept(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/invites/batch-accept" {
			t.Errorf("expected path /v1/invites/batch-accept, got %s", r.URL.Path)
		}

		var body BatchInviteInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body.Invites) != 2 {
			t.Errorf("expected 2 invites, got %d", len(body.Invites))
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Invites.BatchAccept(context.Background(), &BatchInviteInput{
		Invites: []string{"inv_1", "inv_2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvitesBatchDecline(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/invites/batch-decline" {
			t.Errorf("expected path /v1/invites/batch-decline, got %s", r.URL.Path)
		}

		var body BatchInviteInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body.Invites) != 2 {
			t.Errorf("expected 2 invites, got %d", len(body.Invites))
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Invites.BatchDecline(context.Background(), &BatchInviteInput{
		Invites: []string{"inv_1", "inv_2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvitesDeclineAll(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/invites/decline/all" {
			t.Errorf("expected path /v1/invites/decline/all, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Invites.DeclineAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
