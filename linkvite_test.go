package linkvite

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const key = "link_test-api-key"

// testClient creates a test client with a mock server.
func testClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)

	client, err := NewClient(key, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client, server
}

// jsonResponse creates a successful API response.
func jsonResponse(t *testing.T, data any) []byte {
	t.Helper()
	resp := map[string]any{
		"ok":   true,
		"data": data,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}
	return b
}

// jsonError creates an error API response.
func jsonError(t *testing.T, message string) []byte {
	t.Helper()
	resp := map[string]any{
		"ok":    false,
		"error": message,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal error: %v", err)
	}
	return b
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(key)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if client.apiKey != key {
		t.Errorf("expected apiKey to be '%s', got '%s'", key, client.apiKey)
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL to be '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
	if client.httpClient == nil {
		t.Error("expected httpClient to be set")
	}
}

func TestNewClientWithTokens(t *testing.T) {
	token := "test-access-token"
	client, err := NewClientWithTokens(token, "")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if client.accessToken != token {
		t.Errorf("expected accessToken to be '%s', got '%s'", token, client.accessToken)
	}
	if client.authType != AuthTypeAccessToken {
		t.Errorf("expected authType to be AuthTypeAccessToken, got %s", client.authType.String())
	}
}

func TestNewClientWithTokensAndRefreshToken(t *testing.T) {
	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	client, err := NewClientWithTokens(accessToken, refreshToken)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if client.accessToken != accessToken {
		t.Errorf("expected accessToken to be '%s', got '%s'", accessToken, client.accessToken)
	}
	if client.refreshToken != refreshToken {
		t.Errorf("expected refreshToken to be '%s', got '%s'", refreshToken, client.refreshToken)
	}
	if client.authType != AuthTypeAccessToken {
		t.Errorf("expected authType to be AuthTypeAccessToken, got %s", client.authType.String())
	}
}

func TestRefreshAccessToken(t *testing.T) {
	accessToken := "old-access-token"
	refreshToken := "old-refresh-token"
	newAccessToken := "new-access-token"
	newRefreshToken := "new-refresh-token"

	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/auth/token/refresh" {
			var req struct {
				RefreshToken string `json:"refresh_token"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}
			if req.RefreshToken != refreshToken {
				t.Fatalf("expected refresh token '%s', got '%s'", refreshToken, req.RefreshToken)
			}
			_, _ = w.Write(jsonResponse(t, map[string]string{
				"access_token":  newAccessToken,
				"refresh_token": newRefreshToken,
			}))
			return
		}
		t.Fatalf("unexpected request: %s", r.URL.Path)
	})
	defer server.Close()

	client.accessToken = accessToken
	client.refreshToken = refreshToken
	err := client.RefreshAccessToken(context.Background())
	if err != nil {
		t.Fatalf("failed to refresh access token: %v", err)
	}
	if client.accessToken != newAccessToken {
		t.Errorf("expected accessToken to be '%s', got '%s'", newAccessToken, client.accessToken)
	}
	if client.refreshToken != newRefreshToken {
		t.Errorf("expected refreshToken to be '%s', got '%s'", newRefreshToken, client.refreshToken)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	customHTTP := &http.Client{}
	client, err := NewClient(key,
		WithBaseURL("https://custom.api.com"),
		WithHTTPClient(customHTTP),
		WithUserAgent("custom-agent/1.0"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if client.baseURL != "https://custom.api.com" {
		t.Errorf("expected baseURL to be 'https://custom.api.com', got '%s'", client.baseURL)
	}
	if client.httpClient != customHTTP {
		t.Error("expected httpClient to be custom client")
	}
	if client.userAgent != "custom-agent/1.0" {
		t.Errorf("expected userAgent to be 'custom-agent/1.0', got '%s'", client.userAgent)
	}
}

func TestClientAuthorizationHeader(t *testing.T) {
	var receivedAuth string
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("x-api-key")
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.get(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := key
	if receivedAuth != expected {
		t.Errorf("expected x-api-key header '%s', got '%s'", expected, receivedAuth)
	}
}

func TestClientUserAgentHeader(t *testing.T) {
	var receivedUA string
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		_, _ = w.Write(jsonResponse(t, nil))
	})
	defer server.Close()

	err := client.get(context.Background(), "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "linkvite-go/" + Version
	if receivedUA != expected {
		t.Errorf("expected User-Agent header '%s', got '%s'", expected, receivedUA)
	}
}

func TestClientAPIError(t *testing.T) {
	client, server := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(jsonError(t, "resource not found"))
	})
	defer server.Close()

	err := client.get(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "resource not found" {
		t.Errorf("expected message 'resource not found', got '%s'", apiErr.Message)
	}
	if !apiErr.IsNotFound() {
		t.Error("expected IsNotFound() to be true")
	}
}

func TestErrorHelpers(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		check      func(*Error) bool
		expected   bool
	}{
		{"IsNotFound", http.StatusNotFound, (*Error).IsNotFound, true},
		{"IsUnauthorized", http.StatusUnauthorized, (*Error).IsUnauthorized, true},
		{"IsForbidden", http.StatusForbidden, (*Error).IsForbidden, true},
		{"IsRateLimited", http.StatusTooManyRequests, (*Error).IsRateLimited, true},
		{"IsNotFound_false", http.StatusBadRequest, (*Error).IsNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{StatusCode: tt.statusCode}
			if got := tt.check(err); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestPtr(t *testing.T) {
	s := "hello"
	p := Ptr(s)
	if *p != s {
		t.Errorf("expected %s, got %s", s, *p)
	}

	i := 42
	pi := Ptr(i)
	if *pi != i {
		t.Errorf("expected %d, got %d", i, *pi)
	}

	b := true
	pb := Ptr(b)
	if *pb != b {
		t.Errorf("expected %v, got %v", b, *pb)
	}
}

func TestListOptionsToQuery(t *testing.T) {
	tests := []struct {
		name     string
		opts     *ListOptions
		expected string
	}{
		{"nil", nil, ""},
		{"empty", &ListOptions{}, ""},
		{"page only", &ListOptions{Page: 1}, "?page=1"},
		{"limit only", &ListOptions{Limit: 10}, "?limit=10"},
		{"query only", &ListOptions{Query: "test"}, "?q=test"},
		{"multiple", &ListOptions{Page: 2, Limit: 20, Query: "search"}, "?limit=20&page=2&q=search"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.toQuery()
			if got != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}
