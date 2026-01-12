package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRSSFeedsList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/rss-feeds" {
			t.Errorf("expected path /v1/rss-feeds, got %s", r.URL.Path)
		}

		data := []map[string]any{
			{
				"id":    "rss_123",
				"url":   "https://example.com/feed.xml",
				"title": "Example Feed",
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	feeds, err := client.RSSFeeds.List(context.Background(), &ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(feeds) != 1 {
		t.Errorf("expected 1 feed, got %d", len(feeds))
	}
	if feeds[0].ID != "rss_123" {
		t.Errorf("expected ID 'rss_123', got '%s'", feeds[0].ID)
	}
}

func TestRSSFeedsGetItems(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/rss-feeds/rss_123/items" {
			t.Errorf("expected path /v1/rss-feeds/rss_123/items, got %s", r.URL.Path)
		}

		data := []map[string]any{
			{
				"title": "Article 1",
				"link":  "https://example.com/article-1",
			},
			{
				"title": "Article 2",
				"link":  "https://example.com/article-2",
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	items, err := client.RSSFeeds.GetItems(context.Background(), "rss_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if items[0].Title != "Article 1" {
		t.Errorf("expected title 'Article 1', got '%s'", items[0].Title)
	}
}

func TestRSSFeedsCreate(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/rss-feeds" {
			t.Errorf("expected path /v1/rss-feeds, got %s", r.URL.Path)
		}

		var body CreateRSSFeedInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.URL != "https://blog.example.com/rss" {
			t.Errorf("expected url 'https://blog.example.com/rss', got '%s'", body.URL)
		}
		if body.CollectionID != "coll_456" {
			t.Errorf("expected collection_id 'coll_456', got '%s'", body.CollectionID)
		}

		data := map[string]any{
			"id":            "rss_new",
			"url":           "https://blog.example.com/rss",
			"title":         "Example Blog",
			"collection_id": "coll_456",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	feed, err := client.RSSFeeds.Create(context.Background(), &CreateRSSFeedInput{
		URL:          "https://blog.example.com/rss",
		CollectionID: "coll_456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if feed.ID != "rss_new" {
		t.Errorf("expected ID 'rss_new', got '%s'", feed.ID)
	}
}

func TestRSSFeedsEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/rss-feeds/rss_123" {
			t.Errorf("expected path /v1/rss-feeds/rss_123, got %s", r.URL.Path)
		}

		var body EditRSSFeedInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Title != "Updated Title" {
			t.Errorf("expected title 'Updated Title', got '%s'", body.Title)
		}

		data := map[string]any{
			"id":    "rss_123",
			"title": "Updated Title",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	feed, err := client.RSSFeeds.Edit(context.Background(), "rss_123", &EditRSSFeedInput{
		Title: "Updated Title",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if feed.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", feed.Title)
	}
}

func TestRSSFeedsBatchDelete(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/rss-feeds/batch-delete" {
			t.Errorf("expected path /v1/rss-feeds/batch-delete, got %s", r.URL.Path)
		}

		var body BatchDeleteRSSFeedsInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body.Feeds) != 2 {
			t.Errorf("expected 2 feeds, got %d", len(body.Feeds))
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.RSSFeeds.BatchDelete(context.Background(), &BatchDeleteRSSFeedsInput{
		Feeds: []string{"rss_1", "rss_2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
