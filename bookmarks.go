package linkvite

import (
	"context"
	"fmt"
)

// BookmarksService handles bookmark operations.
type BookmarksService struct {
	client *Client
}

// CreateBookmarkInput represents input for creating a bookmark.
type CreateBookmarkInput struct {
	URL         string `json:"url"`
	Collection  string `json:"collection,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Starred     *bool  `json:"starred,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Cover       string `json:"cover,omitempty"`
}

// CreateManualBookmarkInput represents input for creating a bookmark manually.
type CreateManualBookmarkInput struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	Description   string `json:"description,omitempty"`
	Collection    string `json:"collection,omitempty"`
	Cover         string `json:"cover,omitempty"`
	Tags          string `json:"tags,omitempty"`
	Starred       *bool  `json:"starred,omitempty"`
	Favicon       string `json:"favicon,omitempty"`
	AllowComments *bool  `json:"allow_comments,omitempty"`
}

// EditBookmarkInput represents input for editing a bookmark.
type EditBookmarkInput struct {
	Name          string         `json:"name,omitempty"`
	URL           string         `json:"url,omitempty"`
	Tags          string         `json:"tags,omitempty"`
	Status        BookmarkStatus `json:"status,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Collection    *string        `json:"collection,omitempty"`
	AllowComments *bool          `json:"allow_comments,omitempty"`
	Favicon       string         `json:"favicon,omitempty"`
	Type          string         `json:"type,omitempty"`
	Starred       *bool          `json:"starred,omitempty"`
	Cover         string         `json:"cover,omitempty"`
}

// BatchEditBookmarkInput represents input for batch editing bookmarks.
type BatchEditBookmarkInput struct {
	Tags       *string        `json:"tags,omitempty"`
	Status     BookmarkStatus `json:"status,omitempty"`
	Collection string         `json:"collection,omitempty"`
	Type       string         `json:"type,omitempty"`
	Bookmarks  []string       `json:"bookmarks"`
}

// BatchToggleStarInput represents input for batch starring/unstarring.
type BatchToggleStarInput struct {
	Value     *bool    `json:"value"`
	Bookmarks []string `json:"bookmarks"`
}

// BookmarkExistsInput represents input for checking bookmark existence.
type BookmarkExistsInput struct {
	URL        string `json:"url"`
	Collection string `json:"collection,omitempty"`
}

// BookmarkExistsResult represents the result of bookmark existence check.
type BookmarkExistsResult struct {
	Exists   bool      `json:"exists"`
	Bookmark *Bookmark `json:"bookmark"`
}

// ListBookmarksResult represents a paginated list of bookmarks.
type ListBookmarksResult struct {
	Bookmarks  []*BookmarkWithMeta `json:"bookmarks"`
	Pagination *Pagination         `json:"pagination"`
}

// ListBookmarksWithStorageResult represents a paginated list of bookmarks with storage information.
type ListBookmarksWithStorageResult struct {
	Bookmarks  []*BookmarkStorage `json:"bookmarks"`
	Pagination *Pagination        `json:"pagination"`
}

// TabBookmark represents a browser tab.
type TabBookmark struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// CreateTabsBookmarksInput represents input for creating bookmarks from tabs.
type CreateTabsBookmarksInput struct {
	Collection string        `json:"collection,omitempty"`
	Tabs       []TabBookmark `json:"tabs"`
}

// List retrieves bookmarks for the current user.
func (s *BookmarksService) List(ctx context.Context, opts *ListOptions) (*ListBookmarksResult, error) {
	var result ListBookmarksResult
	if err := s.client.get(ctx, "/bookmarks"+opts.toQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAll retrieves all bookmarks without pagination.
func (s *BookmarksService) ListAll(ctx context.Context) (*ListBookmarksResult, error) {
	var result ListBookmarksResult
	if err := s.client.get(ctx, "/bookmarks/all", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a single bookmark by ID.
func (s *BookmarksService) Get(ctx context.Context, id string) (*BookmarkWithMeta, error) {
	var result BookmarkWithMeta
	if err := s.client.get(ctx, "/bookmarks/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetReader retrieves the readable content of a bookmark.
func (s *BookmarksService) GetReader(ctx context.Context, id string) (*BookmarkReader, error) {
	var result BookmarkReader
	if err := s.client.get(ctx, "/bookmarks/"+id+"/reader", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new bookmark.
func (s *BookmarksService) Create(ctx context.Context, input *CreateBookmarkInput) (*Bookmark, error) {
	var result Bookmark
	if err := s.client.post(ctx, "/bookmarks", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateManual creates bookmarks with manual metadata.
func (s *BookmarksService) CreateManual(ctx context.Context, input []*CreateManualBookmarkInput) ([]*Bookmark, error) {
	var result []*Bookmark
	if err := s.client.post(ctx, "/bookmarks/manual", input, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateFromTabs creates bookmarks from browser tabs.
func (s *BookmarksService) CreateFromTabs(ctx context.Context, input *CreateTabsBookmarksInput) ([]*Bookmark, error) {
	var result []*Bookmark
	if err := s.client.post(ctx, "/bookmarks/tabs", input, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Edit updates a bookmark.
func (s *BookmarksService) Edit(ctx context.Context, id string, input *EditBookmarkInput) (*Bookmark, error) {
	var result Bookmark
	if err := s.client.patch(ctx, "/bookmarks/"+id, input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchEdit updates multiple bookmarks at once.
func (s *BookmarksService) BatchEdit(ctx context.Context, input *BatchEditBookmarkInput) error {
	return s.client.patch(ctx, "/bookmarks/batch-edit", input, nil)
}

// Delete deletes a bookmark.
func (s *BookmarksService) Delete(ctx context.Context, id string) error {
	return s.client.delete(ctx, "/bookmarks/"+id)
}

// ToggleStar toggles the starred status of a bookmark.
func (s *BookmarksService) ToggleStar(ctx context.Context, id string) error {
	return s.client.post(ctx, "/bookmarks/"+id+"/star", nil, nil)
}

// BatchToggleStar toggles the starred status of multiple bookmarks.
func (s *BookmarksService) BatchToggleStar(ctx context.Context, input *BatchToggleStarInput) error {
	return s.client.post(ctx, "/bookmarks/batch-star", input, nil)
}

// Exists checks if a bookmark with the given URL exists.
func (s *BookmarksService) Exists(ctx context.Context, input *BookmarkExistsInput) (*BookmarkExistsResult, error) {
	var result BookmarkExistsResult
	if err := s.client.post(ctx, "/bookmarks/exists", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListByTag retrieves bookmarks related to a specific tag.
func (s *BookmarksService) ListByTag(ctx context.Context, tag string, opts *ListOptions) (*ListBookmarksResult, error) {
	path := fmt.Sprintf("/bookmarks/related/tag?tag=%s", tag)
	if opts != nil {
		q := opts.toQuery()
		if q != "" {
			path += "&" + q[1:]
		}
	}
	var result ListBookmarksResult
	if err := s.client.get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListByCollection retrieves bookmarks from a specific collection.
func (s *BookmarksService) ListByCollection(ctx context.Context, collectionID string, opts *ListOptions) (*ListBookmarksResult, error) {
	var result ListBookmarksResult
	if err := s.client.get(ctx, "/bookmarks/from-collection/"+collectionID+opts.toQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RequestArchive requests an archive for a bookmark.
func (s *BookmarksService) RequestArchive(ctx context.Context, id string) error {
	return s.client.post(ctx, "/bookmarks/"+id+"/archive", nil, nil)
}
