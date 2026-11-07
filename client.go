// Package linkvite is the official Go SDK for the Linkvite API.
//
// The services and types are generated from the Linkvite OpenAPI spec, so
// every endpoint, request and response matches the API.
//
// # Authentication
//
// An API key (Settings → Developers) acts as your own account:
//
//	client, err := linkvite.NewClient("link_...")
//
// Apps other people sign in to use OAuth (authorization code + PKCE). The
// idiomatic route is golang.org/x/oauth2, whose HTTP client adds and refreshes
// the token:
//
//	conf := &oauth2.Config{ClientID: "app_...", Endpoint: oauth2.Endpoint{AuthURL: linkvite.AuthURL, TokenURL: linkvite.TokenURL}, ...}
//	client, err := linkvite.NewClientWithTokens("", "", linkvite.WithHTTPClient(conf.Client(ctx, token)))
//
// Or pass a token directly with NewClientWithTokens / WithTokenSource.
//
// # Usage
//
//	b, err := client.Bookmarks.Create(ctx, &linkvite.CreateBookmarkRequest{URL: "https://example.com"})
//	page, err := client.Bookmarks.GetAll(ctx, &linkvite.BookmarksGetAllParams{Q: linkvite.Ptr("go"), Limit: linkvite.Ptr(50)})
//	for b, err := range client.Paginate... // see Paginate
//
// Errors are *Error values: check Status, Errors and RequestID.
package linkvite

//go:generate go run ./internal/gen

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the production API.
	DefaultBaseURL = "https://api.linkvite.io"
	// AuthURL and TokenURL are Linkvite's OAuth 2.1 endpoints (for x/oauth2).
	AuthURL  = DefaultBaseURL + "/v1/oauth/authorize"
	TokenURL = DefaultBaseURL + "/v1/oauth/token"
)

// Client talks to the Linkvite API. Its fields (Bookmarks, Collections, ...)
// are the generated services. Safe for concurrent use.
type Client struct {
	services

	httpClient *http.Client
	baseURL    string
	userAgent  string
	maxRetries int

	apiKey       string
	tokenSource  func(ctx context.Context) (string, error)
	refreshToken string
	onUnauth     func(ctx context.Context) (string, error)
}

// Option configures a Client.
type Option func(*Client) error

// WithBaseURL points the client at another API (e.g. https://api-stg.linkvite.io).
func WithBaseURL(u string) Option {
	return func(c *Client) error {
		u = strings.TrimRight(strings.TrimSpace(u), "/")
		if u == "" {
			return errors.New("linkvite: base URL cannot be empty")
		}
		c.baseURL = u
		return nil
	}
}

// WithHTTPClient sets the HTTP client (timeouts, proxies, or an x/oauth2
// client that adds and refreshes the bearer token).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("linkvite: http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithUserAgent replaces the User-Agent (default linkvite-go/<Version>).
func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(ua) == "" {
			return errors.New("linkvite: user agent cannot be empty")
		}
		c.userAgent = ua
		return nil
	}
}

// WithAPIKey authenticates with an API key (sent as x-api-key).
func WithAPIKey(key string) Option {
	return func(c *Client) error {
		c.apiKey = strings.TrimSpace(key)
		return nil
	}
}

// WithAccessToken authenticates with a fixed bearer token.
func WithAccessToken(token string) Option {
	return func(c *Client) error {
		t := strings.TrimSpace(token)
		c.tokenSource = func(context.Context) (string, error) { return t, nil }
		return nil
	}
}

// WithTokenSource authenticates with the token f returns for each request.
func WithTokenSource(f func(ctx context.Context) (string, error)) Option {
	return func(c *Client) error {
		c.tokenSource = f
		return nil
	}
}

// WithOnUnauthorized is called once when a request gets 401; returning a new
// token retries the request with it.
func WithOnUnauthorized(f func(ctx context.Context) (string, error)) Option {
	return func(c *Client) error {
		c.onUnauth = f
		return nil
	}
}

// WithMaxRetries sets retries for 429, 5xx and network errors on idempotent
// requests (default 2).
func WithMaxRetries(n int) Option {
	return func(c *Client) error {
		c.maxRetries = n
		return nil
	}
}

// NewClient returns a client authenticated with an API key. An empty key is
// allowed for public endpoints (profiles, public collections).
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	return newClient(append([]Option{WithAPIKey(apiKey)}, opts...))
}

// NewClientWithTokens returns a client using a first-party access token; with
// a refresh token, it rotates the session on 401 (RefreshAccessToken).
func NewClientWithTokens(accessToken, refreshToken string, opts ...Option) (*Client, error) {
	c, err := newClient(opts)
	if err != nil {
		return nil, err
	}
	if accessToken != "" && c.tokenSource == nil {
		if err := WithAccessToken(accessToken)(c); err != nil {
			return nil, err
		}
	}
	c.refreshToken = refreshToken
	if refreshToken != "" && c.onUnauth == nil {
		c.onUnauth = func(ctx context.Context) (string, error) {
			if err := c.RefreshAccessToken(ctx); err != nil {
				return "", err
			}
			return c.tokenSource(ctx)
		}
	}
	return c, nil
}

// NewClientWithCredentials signs in with a username/email and password
// (first-party apps; third-party apps should use OAuth or an API key).
func NewClientWithCredentials(ctx context.Context, identifier, password string, opts ...Option) (*Client, error) {
	s, err := Login(ctx, identifier, password, opts...)
	if err != nil {
		return nil, err
	}
	return NewClientWithTokens(s.AccessToken, s.RefreshToken, opts...)
}

func newClient(opts []Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    DefaultBaseURL,
		userAgent:  "linkvite-go/" + Version,
		maxRetries: 2,
	}
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	c.services = newServices(c)
	return c, nil
}

// Session is a first-party sign-in: the user and a token pair.
type Session = AuthResponse

// Login signs in with a username/email and password.
func Login(ctx context.Context, identifier, password string, opts ...Option) (*Session, error) {
	c, err := newClient(opts)
	if err != nil {
		return nil, err
	}
	c.apiKey, c.tokenSource = "", nil
	out := new(Session)
	body := map[string]string{"identifier": identifier, "password": password}
	if err := c.do(ctx, request{method: "POST", path: "/v1/auth/login", body: body}, out); err != nil {
		return nil, err
	}
	return out, nil
}

// RefreshAccessToken rotates a first-party session (NewClientWithTokens).
// Refresh tokens are single-use: the client keeps the new one.
func (c *Client) RefreshAccessToken(ctx context.Context) error {
	if c.refreshToken == "" {
		return errors.New("linkvite: no refresh token")
	}
	out := new(Session)
	body := map[string]string{"refresh_token": c.refreshToken}
	if err := c.do(ctx, request{method: "POST", path: "/v1/auth/token/refresh", body: body, noAuth: true}, out); err != nil {
		return err
	}
	c.refreshToken = out.RefreshToken
	token := out.AccessToken
	c.tokenSource = func(context.Context) (string, error) { return token, nil }
	return nil
}

// Do calls any endpoint, including ones newer than this SDK, decoding the
// response's data into out (nil to discard).
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	return c.do(ctx, request{method: method, path: path, query: query, body: body}, out)
}

/* ----------------------------------------------------------------- transport -- */

type request struct {
	method     string
	path       string
	pathParams map[string]string
	query      url.Values
	body       any
	multipart  bool
	noAuth     bool
}

func (c *Client) url(r request) (string, error) {
	p := r.path
	for k, v := range r.pathParams {
		if v == "" {
			return "", fmt.Errorf("linkvite: missing path parameter %q for %s", k, r.path)
		}
		p = strings.ReplaceAll(p, "{"+k+"}", url.PathEscape(v))
	}
	u := c.baseURL + p
	if len(r.query) > 0 {
		u += "?" + r.query.Encode()
	}
	return u, nil
}

func idempotent(m string) bool {
	switch m {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions:
		return true
	}
	return false
}

func (c *Client) do(ctx context.Context, r request, out any) error {
	u, err := c.url(r)
	if err != nil {
		return err
	}
	var payload []byte
	contentType := ""
	if r.body != nil && !isNilPtr(r.body) {
		if r.multipart {
			payload, contentType, err = encodeMultipart(r.body)
		} else {
			payload, err = json.Marshal(r.body)
			contentType = "application/json"
		}
		if err != nil {
			return fmt.Errorf("linkvite: encode body: %w", err)
		}
	}
	retries := 0
	if idempotent(r.method) {
		retries = c.maxRetries
	}
	refreshed := false
	var override string
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, r.method, u, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		if !r.noAuth {
			if c.apiKey != "" {
				req.Header.Set("x-api-key", c.apiKey)
			}
			token := override
			if token == "" && c.tokenSource != nil {
				if token, err = c.tokenSource(ctx); err != nil {
					return fmt.Errorf("linkvite: token: %w", err)
				}
			}
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if attempt < retries {
				sleep(ctx, backoff(attempt))
				continue
			}
			return &Error{Status: 0, Message: "network error: " + err.Error(), Method: r.method, URL: u, cause: err}
		}
		if res.StatusCode == http.StatusUnauthorized && !refreshed && c.onUnauth != nil && !r.noAuth {
			drain(res)
			refreshed = true
			if t, err := c.onUnauth(ctx); err == nil && t != "" {
				override = t
				attempt--
				continue
			}
			return &Error{Status: 401, Message: "unauthorized", Method: r.method, URL: u}
		}
		if (res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= 500) && attempt < retries {
			wait := retryAfter(res)
			drain(res)
			if wait == 0 {
				wait = backoff(attempt)
			}
			sleep(ctx, wait)
			continue
		}
		return decode(res, r.method, u, out)
	}
}

// drain discards a response we won't read so its connection can be reused.
func drain(res *http.Response) {
	_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 64<<10))
	_ = res.Body.Close()
}

func decode(res *http.Response, method, u string, out any) error {
	defer func() { _ = res.Body.Close() }()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	rid := res.Header.Get("X-Request-ID")
	if !strings.Contains(res.Header.Get("Content-Type"), "json") {
		if res.StatusCode >= 400 {
			return &Error{Status: res.StatusCode, Message: strings.TrimSpace(string(raw)), Method: method, URL: u, RequestID: rid}
		}
		if s, ok := out.(*string); ok {
			*s = string(raw)
		}
		return nil
	}
	var env struct {
		OK        *bool             `json:"ok"`
		Data      json.RawMessage   `json:"data"`
		Message   string            `json:"message"`
		Error     string            `json:"error"`
		Errors    map[string]string `json:"errors"`
		RequestID string            `json:"request_id"`
	}
	_ = json.Unmarshal(raw, &env)
	if res.StatusCode >= 400 || (env.OK != nil && !*env.OK) {
		e := &Error{Status: res.StatusCode, Message: env.Error, Errors: env.Errors, RequestID: env.RequestID, Method: method, URL: u, Body: raw}
		if e.RequestID == "" {
			e.RequestID = rid
		}
		if e.Message == "" {
			e.Message = http.StatusText(res.StatusCode)
		}
		return e
	}
	if out == nil {
		return nil
	}
	data := []byte(env.Data)
	if env.OK == nil {
		data = raw // not enveloped
	}
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("linkvite: decode %s %s: %w", method, u, err)
	}
	return nil
}

func backoff(attempt int) time.Duration {
	d := min(8*time.Second, 400*time.Millisecond<<attempt)
	return time.Duration(float64(d) * (0.75 + rand.Float64()*0.5))
}

func retryAfter(res *http.Response) time.Duration {
	v := res.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if s, err := strconv.Atoi(v); err == nil {
		return min(time.Duration(s)*time.Second, 30*time.Second)
	}
	if t, err := http.ParseTime(v); err == nil {
		return min(time.Until(t), 30*time.Second)
	}
	return 0
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

func isNilPtr(v any) bool {
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Pointer && rv.IsNil()
}

/* ------------------------------------------------------------------ helpers -- */

// Ptr returns a pointer to v, for optional fields and parameters.
func Ptr[T any](v T) *T { return &v }

// File is a file field in a multipart request.
type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
}

func formatParam(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case time.Time:
		return x.Format(time.RFC3339)
	default:
		return fmt.Sprint(x)
	}
}

// encodeMultipart writes the fields of a generated *XRequest (form tags).
func encodeMultipart(body any) ([]byte, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	rv := reflect.Indirect(reflect.ValueOf(body))
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		name := rt.Field(i).Tag.Get("form")
		if name == "" {
			continue
		}
		fv := rv.Field(i)
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}
		if f, ok := fv.Interface().(File); ok {
			if f.Reader == nil {
				continue
			}
			fn := f.Name
			if fn == "" {
				fn = name
			}
			part, err := w.CreateFormFile(name, fn)
			if err != nil {
				return nil, "", err
			}
			if _, err := io.Copy(part, f.Reader); err != nil {
				return nil, "", err
			}
			continue
		}
		if err := w.WriteField(name, formatParam(fv.Interface())); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), w.FormDataContentType(), nil
}
