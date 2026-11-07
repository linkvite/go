package linkvite

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

type recorded struct {
	method, path, query, apiKey, auth, ua, ctype string
	body                                         []byte
	form                                         map[string]string
}

func server(t *testing.T, handler func(w http.ResponseWriter, r *http.Request, n int)) (*httptest.Server, *[]recorded) {
	t.Helper()
	var calls []recorded
	var n int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := recorded{method: r.Method, path: r.URL.EscapedPath(), query: r.URL.RawQuery, apiKey: r.Header.Get("x-api-key"), auth: r.Header.Get("Authorization"), ua: r.Header.Get("User-Agent"), ctype: r.Header.Get("Content-Type")}
		if strings.HasPrefix(rec.ctype, "multipart/") {
			if err := r.ParseMultipartForm(1 << 20); err == nil {
				rec.form = map[string]string{}
				for k, v := range r.MultipartForm.Value {
					rec.form[k] = v[0]
				}
				for k, fh := range r.MultipartForm.File {
					f, _ := fh[0].Open()
					b, _ := io.ReadAll(f)
					rec.form["file:"+k] = fh[0].Filename + "=" + string(b)
				}
			}
		} else {
			rec.body, _ = io.ReadAll(r.Body)
		}
		calls = append(calls, rec)
		handler(w, r, int(atomic.AddInt32(&n, 1)))
	}))
	t.Cleanup(srv.Close)
	return srv, &calls
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func newTest(t *testing.T, srv *httptest.Server, opts ...Option) *Client {
	t.Helper()
	c, err := NewClient("link_test", append([]Option{WithBaseURL(srv.URL)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestRequestBuilding(t *testing.T) {
	srv, calls := server(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 200, map[string]any{"ok": true, "data": map[string]any{"bookmarks": []any{}, "pagination": nil}})
	})
	c := newTest(t, srv)
	_, err := c.Bookmarks.FromCollection(context.Background(), "col 1", &BookmarksFromCollectionParams{Page: Ptr(2), Q: Ptr("go")})
	if err != nil {
		t.Fatal(err)
	}
	got := (*calls)[0]
	if got.method != "GET" || got.path != "/v1/bookmarks/from-collection/col%201" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.query != "page=2&q=go" {
		t.Errorf("query = %q", got.query)
	}
	if got.apiKey != "link_test" || got.auth != "" {
		t.Errorf("auth headers: key=%q auth=%q", got.apiKey, got.auth)
	}
	if !strings.HasPrefix(got.ua, "linkvite-go/") {
		t.Errorf("user agent %q", got.ua)
	}
}

func TestJSONBodyAndEnvelope(t *testing.T) {
	srv, calls := server(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 201, map[string]any{"ok": true, "message": "Bookmark created!", "data": map[string]any{"id": "bmk_1", "url": "https://example.com", "starred": true}})
	})
	c := newTest(t, srv)
	b, err := c.Bookmarks.Create(context.Background(), &CreateBookmarkRequest{URL: "https://example.com", Starred: Ptr(false)})
	if err != nil {
		t.Fatal(err)
	}
	if b.ID != "bmk_1" || !b.Starred {
		t.Errorf("decoded %+v", b)
	}
	var sent map[string]any
	json.Unmarshal((*calls)[0].body, &sent)
	if sent["url"] != "https://example.com" || sent["starred"] != false {
		t.Errorf("sent %s", (*calls)[0].body)
	}
	if _, ok := sent["title"]; ok {
		t.Error("empty optional field was sent")
	}
}

func TestMessageOnlyAndText(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		if strings.HasSuffix(r.URL.Path, "/rss") {
			w.Header().Set("Content-Type", "application/xml")
			io.WriteString(w, "<rss/>")
			return
		}
		writeJSON(w, 200, map[string]any{"ok": true, "message": "Bookmark deleted"})
	})
	c := newTest(t, srv)
	if err := c.Bookmarks.Delete(context.Background(), "bmk_1"); err != nil {
		t.Fatal(err)
	}
	xml, err := c.Profiles.GetCollectionRSS(context.Background(), "kayode", "col_1")
	if err != nil || xml != "<rss/>" {
		t.Errorf("rss = %q, %v", xml, err)
	}
}

func TestMultipart(t *testing.T) {
	srv, calls := server(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 201, map[string]any{"ok": true, "message": "Processing"})
	})
	c := newTest(t, srv)
	err := c.Bookmarks.CreateFile(context.Background(), &BookmarksCreateFileRequest{
		File:    File{Name: "a.txt", Reader: strings.NewReader("hello")},
		Title:   Ptr("hi"),
		Starred: Ptr(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	f := (*calls)[0].form
	if f["title"] != "hi" || f["starred"] != "true" || f["file:file"] != "a.txt=hello" {
		t.Errorf("form = %v", f)
	}
}

func TestErrors(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "Validation error", "errors": map[string]string{"url": "required"}, "request_id": "rid"})
	})
	c := newTest(t, srv)
	_, err := c.Bookmarks.Create(context.Background(), &CreateBookmarkRequest{})
	var e *Error
	if !errors.As(err, &e) || e.Status != 400 || e.Errors["url"] != "required" || e.RequestID != "rid" {
		t.Fatalf("err = %#v", err)
	}
	if !strings.Contains(e.Error(), "rid") {
		t.Errorf("message %q lacks request id", e.Error())
	}
}

func TestRetries(t *testing.T) {
	srv, calls := server(t, func(w http.ResponseWriter, r *http.Request, n int) {
		switch n {
		case 1:
			w.Header().Set("Retry-After", "0")
			writeJSON(w, 429, map[string]any{"ok": false, "error": "slow down"})
		case 2:
			writeJSON(w, 503, map[string]any{"ok": false})
		default:
			writeJSON(w, 200, map[string]any{"ok": true, "data": map[string]any{"id": "u"}})
		}
	})
	c := newTest(t, srv)
	u, err := c.User.Get(context.Background())
	if err != nil || u.ID != "u" || len(*calls) != 3 {
		t.Fatalf("user=%+v err=%v calls=%d", u, err, len(*calls))
	}

	srv2, calls2 := server(t, func(w http.ResponseWriter, r *http.Request, _ int) {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "boom"})
	})
	c2 := newTest(t, srv2)
	if _, err := c2.Bookmarks.Create(context.Background(), &CreateBookmarkRequest{URL: "https://x.y"}); err == nil || len(*calls2) != 1 {
		t.Fatalf("POST retried: calls=%d err=%v", len(*calls2), err)
	}
}

func TestUnauthorizedRefresh(t *testing.T) {
	srv, calls := server(t, func(w http.ResponseWriter, r *http.Request, n int) {
		if r.Header.Get("Authorization") != "Bearer fresh" {
			writeJSON(w, 401, map[string]any{"ok": false, "error": "expired"})
			return
		}
		writeJSON(w, 200, map[string]any{"ok": true, "data": map[string]any{"id": "u"}})
	})
	refreshes := 0
	c, _ := NewClientWithTokens("stale", "", WithBaseURL(srv.URL), WithOnUnauthorized(func(context.Context) (string, error) {
		refreshes++
		return "fresh", nil
	}))
	u, err := c.User.Get(context.Background())
	if err != nil || u.ID != "u" || refreshes != 1 || len(*calls) != 2 {
		t.Fatalf("user=%+v err=%v refreshes=%d calls=%d", u, err, refreshes, len(*calls))
	}
}

func TestPaginate(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request, n int) {
		writeJSON(w, 200, map[string]any{"ok": true, "data": map[string]any{
			"bookmarks":  []any{map[string]any{"id": r.URL.Query().Get("page")}},
			"pagination": map[string]any{"page": n, "last_page": 3, "total": 3, "limit": 1},
		}})
	})
	c := newTest(t, srv)
	var ids []string
	for b, err := range Paginate(context.Background(), func(ctx context.Context, page int) ([]BookmarkResponse, *Pagination, error) {
		r, err := c.Bookmarks.GetAll(ctx, &BookmarksGetAllParams{Page: Ptr(page), Limit: Ptr(1)})
		if err != nil {
			return nil, nil, err
		}
		return r.Bookmarks, r.Pagination, nil
	}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, b.ID)
	}
	if strings.Join(ids, ",") != "1,2,3" {
		t.Errorf("ids = %v", ids)
	}
}

func TestMissingPathParam(t *testing.T) {
	c, _ := NewClient("k")
	if _, err := c.Bookmarks.Get(context.Background(), ""); err == nil || !strings.Contains(err.Error(), `"id"`) {
		t.Errorf("err = %v", err)
	}
}

// Every public operation with an SDK name has a method on its service.
func TestEveryOperationHasAMethod(t *testing.T) {
	raw, err := os.ReadFile("spec/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Paths map[string]map[string]struct {
			Visibility string `json:"x-linkvite-visibility"`
			SDK        string `json:"x-linkvite-sdk"`
		} `json:"paths"`
	}
	json.Unmarshal(raw, &doc)
	c, _ := NewClient("k")
	cv := reflect.ValueOf(c).Elem()
	n := 0
	for _, item := range doc.Paths {
		for _, op := range item {
			if op.Visibility != "public" || op.SDK == "" {
				continue
			}
			ns, m, _ := strings.Cut(op.SDK, ".")
			svc := cv.FieldByName(goNameForTest(ns))
			if !svc.IsValid() || !svc.MethodByName(goNameForTest(m)).IsValid() {
				t.Errorf("no method for %s", op.SDK)
			}
			n++
		}
	}
	if n < 100 {
		t.Errorf("only %d operations", n)
	}
}

func goNameForTest(s string) string {
	for _, ini := range []string{"api", "rss", "csv", "html", "id"} {
		if strings.HasPrefix(strings.ToLower(s), ini) && (len(s) == len(ini) || s[len(ini)] >= 'A' && s[len(ini)] <= 'Z') {
			return strings.ToUpper(ini) + s[len(ini):]
		}
	}
	r := strings.ToUpper(s[:1]) + s[1:]
	for _, ini := range []string{"Csv", "Html", "Rss", "Api", "Id"} {
		r = strings.ReplaceAll(r, ini, strings.ToUpper(ini))
	}
	return r
}
