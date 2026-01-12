package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestUserGet(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user" {
			t.Errorf("expected path /v1/user, got %s", r.URL.Path)
		}

		data := map[string]any{
			"id":           "usr_123",
			"name":         "Test User",
			"email":        "test@example.com",
			"username":     "testuser",
			"account_type": "free",
			"status":       "active",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	user, err := client.User.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "usr_123" {
		t.Errorf("expected ID 'usr_123', got '%s'", user.ID)
	}
	if user.Name != "Test User" {
		t.Errorf("expected name 'Test User', got '%s'", user.Name)
	}
	if user.AccountType != AccountTypeFree {
		t.Errorf("expected account type 'free', got '%s'", user.AccountType)
	}
}

func TestUserEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user" {
			t.Errorf("expected path /v1/user, got %s", r.URL.Path)
		}

		var body EditUserInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Name == nil || *body.Name != "New Name" {
			t.Errorf("expected name 'New Name', got '%v'", body.Name)
		}

		data := map[string]any{
			"id":   "usr_123",
			"name": "New Name",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	user, err := client.User.Edit(context.Background(), &EditUserInput{
		Name: Ptr("New Name"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", user.Name)
	}
}

func TestUserEditSettings(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user/settings" {
			t.Errorf("expected path /v1/user/settings, got %s", r.URL.Path)
		}

		var body EditSettingsInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.PushEnabled == nil || *body.PushEnabled != true {
			t.Errorf("expected push_enabled true, got %v", body.PushEnabled)
		}

		data := map[string]any{
			"id":           "usr_123",
			"push_enabled": true,
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	user, err := client.User.EditSettings(context.Background(), &EditSettingsInput{
		PushEnabled: Ptr(true),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !user.PushEnabled {
		t.Error("expected push_enabled to be true")
	}
}

func TestUserGetStorageUsage(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user/storage" {
			t.Errorf("expected path /v1/user/storage, got %s", r.URL.Path)
		}

		data := map[string]any{
			"used":  1024.5,
			"total": 10240.0,
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	storage, err := client.User.GetStorageUsage(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if storage.Used != 1024.5 {
		t.Errorf("expected used 1024.5, got %f", storage.Used)
	}
	if storage.Total != 10240.0 {
		t.Errorf("expected total 10240.0, got %f", storage.Total)
	}
}

func TestUserGetDetailedStorageUsage(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user/storage/detailed" {
			t.Errorf("expected path /v1/user/storage/detailed, got %s", r.URL.Path)
		}

		data := map[string]any{
			"bookmarks": []*BookmarkStorage{
				{
					ID:          "bkm_123",
					Title:       "Bookmark 1",
					Description: "Description 1",
					Thumbnail:   "Thumbnail 1",
					Icon:        "Icon 1",
					URL:         "https://example.com/1",
					CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					UserID:      "usr_123",
					Size:        1024,
					MimeTypeID:  0,
					Status:      "active",
				},
				{
					ID:          "bkm_124",
					Title:       "Bookmark 2",
					Description: "Description 2",
					Thumbnail:   "Thumbnail 2",
					Icon:        "Icon 2",
					URL:         "https://example.com/2",
					CreatedAt:   time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC),
					UserID:      "usr_123",
					Size:        2048,
					MimeTypeID:  0,
					Status:      "active",
				},
			},
			"pagination": map[string]any{
				"page":  1,
				"limit": 10,
				"total": 2,
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	storage, err := client.User.GetDetailedStorageUsage(context.Background(), &ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(storage.Bookmarks) != 2 {
		t.Errorf("expected 2 bookmarks, got %d", len(storage.Bookmarks))
	}
	if storage.Bookmarks[0].ID != "bkm_123" {
		t.Errorf("expected first bookmark ID 'bkm_123', got '%s'", storage.Bookmarks[0].ID)
	}
	if storage.Bookmarks[1].ID != "bkm_124" {
		t.Errorf("expected second bookmark ID 'bkm_124', got '%s'", storage.Bookmarks[1].ID)
	}
}

func TestUserEmptyTrash(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/user/trash" {
			t.Errorf("expected path /v1/user/trash, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.User.EmptyTrash(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAccountTypeConstants(t *testing.T) {
	if AccountTypeFree != "free" {
		t.Errorf("expected AccountTypeFree to be 'free', got '%s'", AccountTypeFree)
	}
	if AccountTypePro != "pro" {
		t.Errorf("expected AccountTypePro to be 'pro', got '%s'", AccountTypePro)
	}
	if AccountTypeEnterprise != "enterprise" {
		t.Errorf("expected AccountTypeEnterprise to be 'enterprise', got '%s'", AccountTypeEnterprise)
	}
}

func TestUserStatusConstants(t *testing.T) {
	if UserStatusActive != "active" {
		t.Errorf("expected UserStatusActive to be 'active', got '%s'", UserStatusActive)
	}
	if UserStatusSuspended != "suspended" {
		t.Errorf("expected UserStatusSuspended to be 'suspended', got '%s'", UserStatusSuspended)
	}
	if UserStatusDeleted != "deleted" {
		t.Errorf("expected UserStatusDeleted to be 'deleted', got '%s'", UserStatusDeleted)
	}
}
