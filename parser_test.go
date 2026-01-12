package linkvite

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestParserParseLink(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/v1/parser") {
			t.Errorf("expected path /v1/parser, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("url") != "https://example.com/article" {
			t.Errorf("expected url 'https://example.com/article', got '%s'", r.URL.Query().Get("url"))
		}

		data := map[string]any{
			"url":         "https://example.com/article",
			"title":       "Example Article",
			"description": "An interesting article",
			"image":       "https://example.com/image.jpg",
			"author":      "John Doe",
			"site_name":   "Example Site",
		}
		_, _ = w.Write(jsonResponse(t, data))
	})
	defer server.Close()

	parsed, err := client.Parser.ParseLink(context.Background(), "https://example.com/article")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed.URL != "https://example.com/article" {
		t.Errorf("expected url 'https://example.com/article', got '%s'", parsed.URL)
	}
	if parsed.Title != "Example Article" {
		t.Errorf("expected title 'Example Article', got '%s'", parsed.Title)
	}
	if parsed.Description != "An interesting article" {
		t.Errorf("expected description 'An interesting article', got '%s'", parsed.Description)
	}
}

func TestParserImport(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/parser/import" {
			t.Errorf("expected path /v1/parser/import, got %s", r.URL.Path)
		}

		var body ImportParsedItemsInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(body.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(body.Items))
		}
		if body.Collection != "coll_456" {
			t.Errorf("expected collection 'coll_456', got '%s'", body.Collection)
		}

		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.Parser.Import(context.Background(), &ImportParsedItemsInput{
		Items: []ParsedLink{
			{URL: "https://example.com/1", Title: "Article 1"},
			{URL: "https://example.com/2", Title: "Article 2"},
		},
		Collection: "coll_456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
