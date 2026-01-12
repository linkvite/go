package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCommentsList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/comments/bmk_123" {
			t.Errorf("expected path /v1/comments/bmk_123, got %s", r.URL.Path)
		}

		data := []map[string]any{
			{
				"id":          "cmt_123",
				"bookmark_id": "bmk_123",
				"content":     "Great bookmark!",
				"user": map[string]any{
					"id":       "usr_456",
					"name":     "Test User",
					"username": "testuser",
				},
			},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	comments, err := client.Comments.List(context.Background(), "bmk_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].ID != "cmt_123" {
		t.Errorf("expected ID 'cmt_123', got '%s'", comments[0].ID)
	}
	if comments[0].Content != "Great bookmark!" {
		t.Errorf("expected content 'Great bookmark!', got '%s'", comments[0].Content)
	}
}

func TestCommentsCreate(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/comments" {
			t.Errorf("expected path /v1/comments, got %s", r.URL.Path)
		}

		var body CreateCommentInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.BookmarkID != "bmk_123" {
			t.Errorf("expected bookmark_id 'bmk_123', got '%s'", body.BookmarkID)
		}
		if body.Content != "New comment" {
			t.Errorf("expected content 'New comment', got '%s'", body.Content)
		}

		data := map[string]any{
			"id":          "cmt_new",
			"bookmark_id": "bmk_123",
			"content":     "New comment",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	comment, err := client.Comments.Create(context.Background(), &CreateCommentInput{
		BookmarkID: "bmk_123",
		Content:    "New comment",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != "cmt_new" {
		t.Errorf("expected ID 'cmt_new', got '%s'", comment.ID)
	}
}

func TestCommentsDelete(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/comments/cmt_123" {
			t.Errorf("expected path /v1/comments/cmt_123, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Comments.Delete(context.Background(), "cmt_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommentsToggleLike(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/comments/cmt_123/like" {
			t.Errorf("expected path /v1/comments/cmt_123/like, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Comments.ToggleLike(context.Background(), "cmt_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommentsGetLikes(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/comments/cmt_123/likes" {
			t.Errorf("expected path /v1/comments/cmt_123/likes, got %s", r.URL.Path)
		}

		data := []map[string]any{
			{"id": "usr_456", "name": "Liker 1", "username": "liker1"},
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	likes, err := client.Comments.GetLikes(context.Background(), "cmt_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(likes) != 1 {
		t.Errorf("expected 1 like, got %d", len(likes))
	}
	if likes[0].ID != "usr_456" {
		t.Errorf("expected ID 'usr_456', got '%s'", likes[0].ID)
	}
}

func TestCommentsReport(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/comments/cmt_123/report" {
			t.Errorf("expected path /v1/comments/cmt_123/report, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Comments.Report(context.Background(), "cmt_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
