package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestRemindersList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders" {
			t.Errorf("expected path /v1/reminders, got %s", r.URL.Path)
		}

		data := map[string]any{
			"reminders": []map[string]any{
				{
					"id":          "rem_123",
					"bookmark_id": "bmk_456",
					"remind_at":   "2025-01-15T10:00:00Z",
					"completed":   false,
				},
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

	result, err := client.Reminders.List(context.Background(), &ListOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Reminders) != 1 {
		t.Errorf("expected 1 reminder, got %d", len(result.Reminders))
	}
	if result.Reminders[0].ID != "rem_123" {
		t.Errorf("expected ID 'rem_123', got '%s'", result.Reminders[0].ID)
	}
}

func TestRemindersGet(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders/rem_123" {
			t.Errorf("expected path /v1/reminders/rem_123, got %s", r.URL.Path)
		}

		data := map[string]any{
			"id":          "rem_123",
			"bookmark_id": "bmk_456",
			"remind_at":   "2025-01-15T10:00:00Z",
			"completed":   false,
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	reminder, err := client.Reminders.Get(context.Background(), "rem_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reminder.ID != "rem_123" {
		t.Errorf("expected ID 'rem_123', got '%s'", reminder.ID)
	}
}

func TestRemindersCreate(t *testing.T) {
	remindAt := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders" {
			t.Errorf("expected path /v1/reminders, got %s", r.URL.Path)
		}

		var body CreateReminderInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.BookmarkID != "bmk_456" {
			t.Errorf("expected bookmark_id 'bmk_456', got '%s'", body.BookmarkID)
		}
		if body.Note != "Review this article" {
			t.Errorf("expected note 'Review this article', got '%s'", body.Note)
		}

		data := map[string]any{
			"id":          "rem_new",
			"bookmark_id": "bmk_456",
			"remind_at":   "2025-01-15T10:00:00Z",
			"note":        "Review this article",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	reminder, err := client.Reminders.Create(context.Background(), &CreateReminderInput{
		BookmarkID: "bmk_456",
		RemindAt:   remindAt,
		Note:       "Review this article",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reminder.ID != "rem_new" {
		t.Errorf("expected ID 'rem_new', got '%s'", reminder.ID)
	}
}

func TestRemindersEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders/rem_123" {
			t.Errorf("expected path /v1/reminders/rem_123, got %s", r.URL.Path)
		}

		var body EditReminderInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Completed == nil || *body.Completed != true {
			t.Errorf("expected completed true, got %v", body.Completed)
		}

		data := map[string]any{
			"id":        "rem_123",
			"completed": true,
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	reminder, err := client.Reminders.Edit(context.Background(), "rem_123", &EditReminderInput{
		Completed: Ptr(true),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reminder.Completed {
		t.Error("expected completed to be true")
	}
}

func TestRemindersDelete(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders/rem_123" {
			t.Errorf("expected path /v1/reminders/rem_123, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Reminders.Delete(context.Background(), "rem_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemindersBatchEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders/batch-edit" {
			t.Errorf("expected path /v1/reminders/batch-edit, got %s", r.URL.Path)
		}

		var body BatchEditRemindersInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body.Reminders) != 2 {
			t.Errorf("expected 2 reminders, got %d", len(body.Reminders))
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Reminders.BatchEdit(context.Background(), &BatchEditRemindersInput{
		Reminders: []string{"rem_1", "rem_2"},
		Completed: Ptr(true),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemindersBatchDelete(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders/batch-delete" {
			t.Errorf("expected path /v1/reminders/batch-delete, got %s", r.URL.Path)
		}

		var body BatchDeleteRemindersInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body.Reminders) != 2 {
			t.Errorf("expected 2 reminders, got %d", len(body.Reminders))
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Reminders.BatchDelete(context.Background(), &BatchDeleteRemindersInput{
		Reminders: []string{"rem_1", "rem_2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemindersClear(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/reminders/clear" {
			t.Errorf("expected path /v1/reminders/clear, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Reminders.Clear(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
