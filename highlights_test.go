package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestHighlightsList(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/highlights/bmk_123" {
			t.Errorf("expected path /v1/highlights/bmk_123, got %s", r.URL.Path)
		}

		data := map[string]any{
			"highlights": []map[string]any{
				{
					"id":          "hil_123",
					"bookmark_id": "bmk_123",
					"exact":       "highlighted text",
					"start":       0,
					"end":         16,
					"color":       "lv-highlighter-1",
				},
			},
			"bookmark": "bmk_123",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	result, err := client.Highlights.List(context.Background(), "bmk_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Highlights) != 1 {
		t.Errorf("expected 1 highlight, got %d", len(result.Highlights))
	}
	if result.Highlights[0].ID != "hil_123" {
		t.Errorf("expected ID 'hil_123', got '%s'", result.Highlights[0].ID)
	}
	if result.Highlights[0].Exact != "highlighted text" {
		t.Errorf("expected exact 'highlighted text', got '%s'", result.Highlights[0].Exact)
	}
}

func TestHighlightsCreate(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/highlights" {
			t.Errorf("expected path /v1/highlights, got %s", r.URL.Path)
		}

		var body CreateHighlightInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.ID != "hil_new" {
			t.Errorf("expected ID 'hil_new', got '%s'", body.ID)
		}
		if body.BookmarkID != "bmk_123" {
			t.Errorf("expected bookmark_id 'bmk_123', got '%s'", body.BookmarkID)
		}
		if body.Color != HighlightColor1 {
			t.Errorf("expected color 'lv-highlighter-1', got '%s'", body.Color)
		}
		if body.Start == nil || *body.Start != 0 {
			t.Errorf("expected start 0, got %v", body.Start)
		}
		if body.End == nil || *body.End != 50 {
			t.Errorf("expected end 50, got %v", body.End)
		}
		if body.Exact != "test text" {
			t.Errorf("expected exact 'test text', got '%s'", body.Exact)
		}

		data := map[string]any{
			"highlight": map[string]any{
				"id":          "hil_new",
				"bookmark_id": "bmk_123",
				"type":        "Annotation",
			},
			"previous_id": "",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	result, err := client.Highlights.Create(context.Background(), &CreateHighlightInput{
		ID:         "hil_new",
		BookmarkID: "bmk_123",
		Color:      HighlightColor1,
		Start:      Ptr(0),
		End:        Ptr(50),
		Exact:      "test text",
		Note:       "my note",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Highlight.ID != "hil_new" {
		t.Errorf("expected ID 'hil_new', got '%s'", result.Highlight.ID)
	}
}

func TestHighlightsEdit(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/v1/highlights/hil_123" {
			t.Errorf("expected path /v1/highlights/hil_123, got %s", r.URL.Path)
		}

		var body EditHighlightInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if body.Note != "updated note" {
			t.Errorf("expected note 'updated note', got '%s'", body.Note)
		}

		data := map[string]any{
			"id":   "hil_123",
			"note": "updated note",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	highlight, err := client.Highlights.Edit(context.Background(), "hil_123", &EditHighlightInput{
		Note:  "updated note",
		Color: HighlightColor2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if highlight.Note != "updated note" {
		t.Errorf("expected note 'updated note', got '%s'", highlight.Note)
	}
}

func TestHighlightsDelete(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/highlights/hil_123" {
			t.Errorf("expected path /v1/highlights/hil_123, got %s", r.URL.Path)
		}
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Highlights.Delete(context.Background(), "hil_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHighlightColorConstants(t *testing.T) {
	tests := []struct {
		color    HighlightColor
		expected string
	}{
		{HighlightColor1, "lv-highlighter-1"},
		{HighlightColor2, "lv-highlighter-2"},
		{HighlightColor3, "lv-highlighter-3"},
		{HighlightColor4, "lv-highlighter-4"},
		{HighlightColor5, "lv-highlighter-5"},
	}

	for _, tt := range tests {
		if string(tt.color) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.color)
		}
	}
}
