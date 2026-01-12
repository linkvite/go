package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestBookmarksList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/v1/bookmarks") {
			t.Errorf("expected path /v1/bookmarks, got %s", r.URL.Path)
		}

		data := map[string]any{
			"bookmarks": []map[string]any{
				{"id": "bmk_123", "title": "Test Bookmark", "url": "https://example.com"},
			},
			"pagination": map[string]any{
				"total":     1,
				"page":      1,
				"limit":     20,
				"last_page": 1,
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	result, err := client.Bookmarks.List(context.Background(), &ListOptions{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Bookmarks) != 1 {
		t.Errorf("expected 1 bookmark, got %d", len(result.Bookmarks))
	}
	if result.Bookmarks[0].ID != "bmk_123" {
		t.Errorf("expected ID 'bmk_123', got '%s'", result.Bookmarks[0].ID)
	}
}

func TestBookmarksGet(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/bookmarks/bmk_123" {
			t.Errorf("expected path /v1/bookmarks/bmk_123, got %s", r.URL.Path)
		}

		data := map[string]any{
			"id":    "bmk_123",
			"title": "Test Bookmark",
			"url":   "https://example.com",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	bookmark, err := client.Bookmarks.Get(context.Background(), "bmk_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bookmark.ID != "bmk_123" {
		t.Errorf("expected ID 'bmk_123', got '%s'", bookmark.ID)
	}
	if bookmark.Title != "Test Bookmark" {
		t.Errorf("expected title 'Test Bookmark', got '%s'", bookmark.Title)
	}
}

func TestBookmarksCreate(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/bookmarks" {
			t.Errorf("expected path /v1/bookmarks, got %s", r.URL.Path)
		}

		var body CreateBookmarkInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.URL != "https://example.com" {
			t.Errorf("expected URL 'https://example.com', got '%s'", body.URL)
		}

		data := map[string]any{
			"id":    "bmk_new",
			"title": "New Bookmark",
			"url":   body.URL,
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	bookmark, err := client.Bookmarks.Create(context.Background(), &CreateBookmarkInput{
		URL:   "https://example.com",
		Title: "New Bookmark",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bookmark.ID != "bmk_new" {
		t.Errorf("expected ID 'bmk_new', got '%s'", bookmark.ID)
	}
}

func TestBookmarksEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/bookmarks/bmk_123" {
			t.Errorf("expected path /v1/bookmarks/bmk_123, got %s", r.URL.Path)
		}

		data := map[string]any{
			"id":    "bmk_123",
			"title": "Updated Title",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	bookmark, err := client.Bookmarks.Edit(context.Background(), "bmk_123", &EditBookmarkInput{
		Name: "Updated Title",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bookmark.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", bookmark.Title)
	}
}

func TestBookmarksDelete(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/bookmarks/bmk_123" {
			t.Errorf("expected path /v1/bookmarks/bmk_123, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Bookmarks.Delete(context.Background(), "bmk_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookmarksToggleStar(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/bookmarks/bmk_123/star" {
			t.Errorf("expected path /v1/bookmarks/bmk_123/star, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Bookmarks.ToggleStar(context.Background(), "bmk_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBookmarksExists(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/bookmarks/exists" {
			t.Errorf("expected path /v1/bookmarks/exists, got %s", r.URL.Path)
		}

		data := map[string]any{
			"exists": true,
			"bookmark": map[string]any{
				"id":  "bmk_123",
				"url": "https://example.com",
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	result, err := client.Bookmarks.Exists(context.Background(), &BookmarkExistsInput{
		URL: "https://example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Exists {
		t.Error("expected Exists to be true")
	}
	if result.Bookmark == nil {
		t.Error("expected Bookmark to be set")
	}
}

func TestBookmarkType(t *testing.T) {
	now := time.Now()
	bookmark := Bookmark{
		ID:            "bmk_123",
		UserID:        "usr_456",
		Title:         "Test",
		URL:           "https://example.com",
		Status:        BookmarkStatusActive,
		ArchiveStatus: BookmarkArchiveStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if bookmark.Status != BookmarkStatusActive {
		t.Errorf("expected status Active, got %s", bookmark.Status)
	}
	if bookmark.ArchiveStatus != BookmarkArchiveStatusPending {
		t.Errorf("expected archive status Pending, got %s", bookmark.ArchiveStatus)
	}
}
