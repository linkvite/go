package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAPIKeysList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api-keys" {
			t.Errorf("expected path /v1/api-keys, got %s", r.URL.Path)
		}

		data := []map[string]any{
			{
				"id":     "key_123",
				"name":   "Production Key",
				"scopes": []string{"read", "write"},
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	keys, err := client.APIKeys.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
	if keys[0].ID != "key_123" {
		t.Errorf("expected ID 'key_123', got '%s'", keys[0].ID)
	}
	if keys[0].Name != "Production Key" {
		t.Errorf("expected name 'Production Key', got '%s'", keys[0].Name)
	}
}

func TestAPIKeysCreate(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api-keys" {
			t.Errorf("expected path /v1/api-keys, got %s", r.URL.Path)
		}

		var body CreateAPIKeyInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Name != "New API Key" {
			t.Errorf("expected name 'New API Key', got '%s'", body.Name)
		}
		if len(body.Scopes) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(body.Scopes))
		}

		data := map[string]any{
			"id":     "key_new",
			"name":   "New API Key",
			"scopes": []string{"read", "write"},
			"key":    "link_sk_live_abcdef123456",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	result, err := client.APIKeys.Create(context.Background(), &CreateAPIKeyInput{
		Name:   "New API Key",
		Scopes: []string{"read", "write"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "key_new" {
		t.Errorf("expected ID 'key_new', got '%s'", result.ID)
	}
	if result.Key != "link_sk_live_abcdef123456" {
		t.Errorf("expected key 'link_sk_live_abcdef123456', got '%s'", result.Key)
	}
}

func TestAPIKeysEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/api-keys/key_123" {
			t.Errorf("expected path /v1/api-keys/key_123, got %s", r.URL.Path)
		}

		var body EditAPIKeyInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Name != "Renamed Key" {
			t.Errorf("expected name 'Renamed Key', got '%s'", body.Name)
		}

		data := map[string]any{
			"id":   "key_123",
			"name": "Renamed Key",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	key, err := client.APIKeys.Edit(context.Background(), "key_123", &EditAPIKeyInput{
		Name: "Renamed Key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.Name != "Renamed Key" {
		t.Errorf("expected name 'Renamed Key', got '%s'", key.Name)
	}
}
