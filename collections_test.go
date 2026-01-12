package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCollectionsList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections" {
			t.Errorf("expected path /v1/collections, got %s", r.URL.Path)
		}

		data := []map[string]any{
			{"id": "coll_123", "name": "Test Collection"},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	collections, err := client.Collections.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collections) != 1 {
		t.Errorf("expected 1 collection, got %d", len(collections))
	}
	if collections[0].ID != "coll_123" {
		t.Errorf("expected ID 'coll_123', got '%s'", collections[0].ID)
	}
}

func TestCollectionsGet(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123" {
			t.Errorf("expected path /v1/collections/coll_123, got %s", r.URL.Path)
		}

		data := map[string]any{
			"id":   "coll_123",
			"name": "Test Collection",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	collection, err := client.Collections.Get(context.Background(), "coll_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if collection.ID != "coll_123" {
		t.Errorf("expected ID 'coll_123', got '%s'", collection.ID)
	}
}

func TestCollectionsCreate(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections" {
			t.Errorf("expected path /v1/collections, got %s", r.URL.Path)
		}

		var body CreateCollectionInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Name != "New Collection" {
			t.Errorf("expected name 'New Collection', got '%s'", body.Name)
		}

		data := map[string]any{
			"id":   "coll_new",
			"name": body.Name,
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	collection, err := client.Collections.Create(context.Background(), &CreateCollectionInput{
		Name:        "New Collection",
		Description: "A test collection",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if collection.ID != "coll_new" {
		t.Errorf("expected ID 'coll_new', got '%s'", collection.ID)
	}
}

func TestCollectionsEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123" {
			t.Errorf("expected path /v1/collections/coll_123, got %s", r.URL.Path)
		}

		data := map[string]any{
			"id":   "coll_123",
			"name": "Updated Name",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	collection, err := client.Collections.Edit(context.Background(), "coll_123", &EditCollectionInput{
		Name: "Updated Name",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if collection.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", collection.Name)
	}
}

func TestCollectionsDelete(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123" {
			t.Errorf("expected path /v1/collections/coll_123, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Collections.Delete(context.Background(), "coll_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCollectionsGetMembers(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123/members" {
			t.Errorf("expected path /v1/collections/coll_123/members, got %s", r.URL.Path)
		}

		data := map[string]any{
			"role":     "admin",
			"owner_id": "usr_456",
			"members": []map[string]any{
				{"id": "usr_789", "name": "Test User", "role": "viewer"},
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	result, err := client.Collections.GetMembers(context.Background(), "coll_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Role != RoleLevelAdmin {
		t.Errorf("expected role 'admin', got '%s'", result.Role)
	}
	if len(result.Members) != 1 {
		t.Errorf("expected 1 member, got %d", len(result.Members))
	}
}

func TestCollectionsInviteUser(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123/send-invite" {
			t.Errorf("expected path /v1/collections/coll_123/send-invite, got %s", r.URL.Path)
		}

		var body InviteUserInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Role != RoleLevelViewer {
			t.Errorf("expected role 'viewer', got '%s'", body.Role)
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Collections.InviteUser(context.Background(), "coll_123", &InviteUserInput{
		Identifier: "user@example.com",
		Role:       RoleLevelViewer,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCollectionsLeave(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123/leave" {
			t.Errorf("expected path /v1/collections/coll_123/leave, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Collections.Leave(context.Background(), "coll_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoleLevelConstants(t *testing.T) {
	if RoleLevelAdmin != "admin" {
		t.Errorf("expected RoleLevelAdmin to be 'admin', got '%s'", RoleLevelAdmin)
	}
	if RoleLevelViewer != "viewer" {
		t.Errorf("expected RoleLevelViewer to be 'viewer', got '%s'", RoleLevelViewer)
	}
}

func TestCollectionsMove(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123/move" {
			t.Errorf("expected path /v1/collections/coll_123/move, got %s", r.URL.Path)
		}

		var body MoveCollectionInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if *body.Parent != "coll_456" {
			t.Errorf("expected parent 'coll_456', got '%s'", *body.Parent)
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	parentID := "coll_456"
	err := client.Collections.Move(context.Background(), "coll_123", &MoveCollectionInput{
		Parent: &parentID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCollectionsGetPresence(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/collections/coll_123/presence" {
			t.Errorf("expected path /v1/collections/coll_123/presence, got %s", r.URL.Path)
		}

		data := map[string]any{
			"online_clients": 2,
			"online_members": 1,
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	result, err := client.Collections.GetPresence(context.Background(), "coll_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OnlineClients != 2 {
		t.Errorf("expected 2 online clients, got %d", result.OnlineClients)
	}
	if result.OnlineMembers != 1 {
		t.Errorf("expected 1 online member, got %d", result.OnlineMembers)
	}
}
