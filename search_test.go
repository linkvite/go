package linkvite

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestSearchSearch(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/v1/search") {
			t.Errorf("expected path /v1/search, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "golang" {
			t.Errorf("expected query 'golang', got '%s'", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("type") != "bookmarks" {
			t.Errorf("expected type 'bookmarks', got '%s'", r.URL.Query().Get("type"))
		}

		data := map[string]any{
			"bookmarks": []map[string]any{
				{"id": "bmk_123", "title": "Go Programming", "url": "https://golang.org"},
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

	result, err := client.Search.Search(context.Background(), &SearchInput{
		Query: "golang",
		Type:  "bookmarks",
		Page:  1,
		Limit: 20,
	})
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

func TestSearchSearchCollections(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		data := map[string]any{
			"collections": []map[string]any{
				{"id": "coll_123", "name": "Dev Resources"},
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

	result, err := client.Search.Search(context.Background(), &SearchInput{
		Query: "dev",
		Type:  "collections",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Collections) != 1 {
		t.Errorf("expected 1 collection, got %d", len(result.Collections))
	}
}

func TestSearchSearchUsers(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		data := map[string]any{
			"users": []map[string]any{
				{"id": "usr_123", "name": "John Doe", "username": "johndoe"},
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

	result, err := client.Search.Search(context.Background(), &SearchInput{
		Query: "john",
		Type:  "users",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Users) != 1 {
		t.Errorf("expected 1 user, got %d", len(result.Users))
	}
}

func TestSearchUnsplash(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/v1/search/unsplash") {
			t.Errorf("expected path /v1/search/unsplash, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "nature" {
			t.Errorf("expected query 'nature', got '%s'", r.URL.Query().Get("q"))
		}

		data := []map[string]any{
			{
				"id":          "photo_123",
				"description": "Beautiful mountain",
				"urls": map[string]string{
					"raw":     "https://unsplash.com/raw",
					"full":    "https://unsplash.com/full",
					"regular": "https://unsplash.com/regular",
					"small":   "https://unsplash.com/small",
					"thumb":   "https://unsplash.com/thumb",
				},
				"user": map[string]string{
					"name":     "Photographer",
					"username": "photog",
				},
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	photos, err := client.Search.SearchUnsplash(context.Background(), "nature", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(photos) != 1 {
		t.Errorf("expected 1 photo, got %d", len(photos))
	}
	if photos[0].ID != "photo_123" {
		t.Errorf("expected ID 'photo_123', got '%s'", photos[0].ID)
	}
}

func TestSearchTrackUnsplashDownload(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/search/unsplash/track" {
			t.Errorf("expected path /v1/search/unsplash/track, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Search.TrackUnsplashDownload(context.Background(), "photo_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchTenor(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/v1/search/tenor") {
			t.Errorf("expected path /v1/search/tenor, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "happy" {
			t.Errorf("expected query 'happy', got '%s'", r.URL.Query().Get("q"))
		}

		data := []map[string]any{
			{
				"id":    "gif_123",
				"title": "Happy Dance",
				"media": []map[string]any{
					{
						"gif":     map[string]string{"url": "https://tenor.com/gif.gif"},
						"tinygif": map[string]string{"url": "https://tenor.com/tiny.gif"},
					},
				},
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	gifs, err := client.Search.SearchTenor(context.Background(), "happy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gifs) != 1 {
		t.Errorf("expected 1 gif, got %d", len(gifs))
	}
	if gifs[0].ID != "gif_123" {
		t.Errorf("expected ID 'gif_123', got '%s'", gifs[0].ID)
	}
}
