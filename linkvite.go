// Package linkvite provides a Go SDK for the Linkvite API.
//
// Linkvite is a bookmarking service that allows users to save, organize,
// and share bookmarks. This SDK provides a convenient way to interact
// with the Linkvite API from Go applications.
//
// Basic usage:
//
//	client, err := linkvite.NewClient("link_your-api-key")
//
//	if err != nil {
//		log.Fatalf("failed to create client: %v", err)
//	}
//
//	// Get current user
//	user, err := client.User.Get(ctx)
//
//	// Create a bookmark
//	bookmark, err := client.Bookmarks.Create(ctx, &linkvite.CreateBookmarkInput{
//		URL: "https://example.com",
//	})
package linkvite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default Linkvite API base URL.
	DefaultBaseURL = "https://api.linkvite.io"
)

// AuthType represents the type of authentication used by the client.
type AuthType int

const (
	// AuthTypeUnset indicates that no authentication method has been set.
	AuthTypeUnset AuthType = iota

	// AuthTypeAPIKey indicates that the client is using API key authentication.
	AuthTypeAPIKey

	// AuthTypeAccessToken indicates that the client is using access token authentication.
	AuthTypeAccessToken
)

func (a AuthType) String() string {
	switch a {
	case AuthTypeAPIKey:
		return "api_key"
	case AuthTypeAccessToken:
		return "access_token"
	default:
		return "unset"
	}
}

// Client is the Linkvite API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string

	authType     AuthType
	apiKey       string
	accessToken  string
	refreshToken string

	User        *UserService
	Bookmarks   *BookmarksService
	Collections *CollectionsService
	Comments    *CommentsService
	Highlights  *HighlightsService
	Reminders   *RemindersService
	Invites     *InvitesService
	Search      *SearchService
	RSSFeeds    *RSSFeedsService
	APIKeys     *APIKeysService
	Parser      *ParserService
}

// Option configures a Client.
type Option func(*Client) error

// WithBaseURL sets a custom base URL.
func WithBaseURL(u string) Option {
	return func(c *Client) error {
		ul := strings.TrimSpace(u)
		if ul == "" {
			return fmt.Errorf("linkvite: baseURL cannot be empty")
		}
		c.baseURL = ul
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return fmt.Errorf("linkvite: http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		u := strings.TrimSpace(ua)
		if u == "" {
			return fmt.Errorf("linkvite: userAgent cannot be empty")
		}
		c.userAgent = u
		return nil
	}
}

// WithAPIKey configures API key authentication.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		k := strings.TrimSpace(apiKey)
		if k == "" {
			return fmt.Errorf("linkvite: apiKey cannot be empty")
		}
		if !strings.HasPrefix(k, "link_") {
			return fmt.Errorf("linkvite: invalid apiKey format")
		}

		// conflict detection: do not allow both auth methods
		if c.authType != AuthTypeUnset && c.authType != AuthTypeAPIKey {
			return fmt.Errorf("linkvite: cannot set apiKey when authType is %s", c.authType.String())
		}

		c.apiKey = k
		c.authType = AuthTypeAPIKey
		return nil
	}
}

// WithAccessToken configures access-token authentication.
func WithAccessToken(token string) Option {
	return func(c *Client) error {
		t := strings.TrimSpace(token)
		if t == "" {
			return fmt.Errorf("linkvite: accessToken cannot be empty")
		}

		// conflict detection: do not allow both auth methods
		if c.authType != AuthTypeUnset && c.authType != AuthTypeAPIKey {
			return fmt.Errorf("linkvite: cannot set accessToken when authType is %s", c.authType.String())
		}

		c.accessToken = t
		c.authType = AuthTypeAccessToken
		return nil
	}
}

// WithRefreshToken sets the refresh token for token refresh operations.
func WithRefreshToken(token string) Option {
	return func(c *Client) error {
		t := strings.TrimSpace(token)
		if t == "" {
			return fmt.Errorf("linkvite: refreshToken cannot be empty")
		}
		c.refreshToken = t
		return nil
	}
}

// WithTokens configures access-token authentication and adds a refresh token if provided.
func WithTokens(accessToken, refreshToken string) Option {
	return func(c *Client) error {
		if err := WithAccessToken(accessToken)(c); err != nil {
			return err
		}
		if refreshToken != "" {
			if err := WithRefreshToken(refreshToken)(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// NewClient creates a new Linkvite API client.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    DefaultBaseURL,
		userAgent:  fmt.Sprintf("linkvite-go/%s", Version),
		authType:   AuthTypeUnset,
	}

	if strings.TrimSpace(apiKey) != "" {
		if err := WithAPIKey(apiKey)(c); err != nil {
			return nil, err
		}
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	switch c.authType {
	case AuthTypeUnset:
		return nil, fmt.Errorf("linkvite: either apiKey or accessToken is required")
	case AuthTypeAPIKey:
		if strings.TrimSpace(c.apiKey) == "" {
			return nil, fmt.Errorf("linkvite: apiKey is required for AuthTypeAPIKey")
		}
	case AuthTypeAccessToken:
		if strings.TrimSpace(c.accessToken) == "" {
			return nil, fmt.Errorf("linkvite: accessToken is required for AuthTypeAccessToken")
		}
	default:
		return nil, fmt.Errorf("linkvite: invalid authType %d", int(c.authType))
	}

	c.User = &UserService{client: c}
	c.Bookmarks = &BookmarksService{client: c}
	c.Collections = &CollectionsService{client: c}
	c.Comments = &CommentsService{client: c}
	c.Highlights = &HighlightsService{client: c}
	c.Reminders = &RemindersService{client: c}
	c.Invites = &InvitesService{client: c}
	c.Search = &SearchService{client: c}
	c.RSSFeeds = &RSSFeedsService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.Parser = &ParserService{client: c}

	return c, nil
}

// NewClientWithTokens creates a new Linkvite API client using access and refresh tokens.
func NewClientWithTokens(accessToken, refreshToken string, opts ...Option) (*Client, error) {
	opts = append([]Option{WithTokens(accessToken, refreshToken)}, opts...)
	return NewClient("", opts...)
}

// RefreshAccessToken refreshes the access token using the refresh token
func (c *Client) RefreshAccessToken(ctx context.Context) error {
	if c.refreshToken == "" {
		return fmt.Errorf("linkvite: no refresh token available")
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	err := c.post(ctx, "/auth/token/refresh", map[string]string{"refresh_token": c.refreshToken}, &result)
	if err != nil {
		return err
	}
	c.accessToken = result.AccessToken
	c.refreshToken = result.RefreshToken
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
	reqURL, err := url.Parse(c.baseURL + "/v1" + path)
	if err != nil {
		return fmt.Errorf("linkvite: invalid URL: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("linkvite: marshal error: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("linkvite: request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	switch c.authType {
	case AuthTypeAPIKey:
		req.Header.Set("x-api-key", c.apiKey)
	case AuthTypeAccessToken:
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("linkvite: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("linkvite: read error: %w", err)
	}

	var apiResp struct {
		OK      bool   `json:"ok"`
		Data    any    `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
		Errors  any    `json:"errors,omitempty"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("linkvite: parse error: %w", err)
	}

	if !apiResp.OK {
		return &Error{
			StatusCode: resp.StatusCode,
			Message:    apiResp.Error,
			Errors:     apiResp.Errors,
		}
	}

	if result != nil && apiResp.Data != nil {
		dataBytes, _ := json.Marshal(apiResp.Data)
		if err := json.Unmarshal(dataBytes, result); err != nil {
			return fmt.Errorf("linkvite: unmarshal error: %w", err)
		}
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) post(ctx context.Context, path string, body, result any) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

func (c *Client) patch(ctx context.Context, path string, body, result any) error {
	return c.do(ctx, http.MethodPatch, path, body, result)
}

func (c *Client) delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// Error represents an API error.
type Error struct {
	StatusCode int
	Message    string
	Errors     any
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("linkvite: %s", e.Message)
	}
	return fmt.Sprintf("linkvite: API error (status %d)", e.StatusCode)
}

// IsNotFound returns true if this is a 404 error.
func (e *Error) IsNotFound() bool { return e.StatusCode == http.StatusNotFound }

// IsUnauthorized returns true if this is a 401 error.
func (e *Error) IsUnauthorized() bool { return e.StatusCode == http.StatusUnauthorized }

// IsForbidden returns true if this is a 403 error.
func (e *Error) IsForbidden() bool { return e.StatusCode == http.StatusForbidden }

// IsRateLimited returns true if this is a 429 error.
func (e *Error) IsRateLimited() bool { return e.StatusCode == http.StatusTooManyRequests }

// Ptr returns a pointer to the value.
func Ptr[T any](v T) *T { return &v }
