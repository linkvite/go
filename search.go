package linkvite

import (
	"context"
	"fmt"
	"net/url"
)

// SearchService handles search operations.
type SearchService struct {
	client *Client
}

// SearchInput represents input for searching.
type SearchInput struct {
	Query string
	Type  string // "bookmarks", "collections", "users"
	Page  int64
	Limit int64
}

// SearchResult represents a search result.
type SearchResult struct {
	Bookmarks   []*BookmarkWithMeta   `json:"bookmarks,omitempty"`
	Collections []*CollectionWithMeta `json:"collections,omitempty"`
	Users       []*PublicProfile      `json:"users,omitempty"`
	Pagination  *Pagination           `json:"pagination"`
}

// UnsplashPhoto represents a photo from Unsplash.
type UnsplashPhoto struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	URLs        struct {
		Raw     string `json:"raw"`
		Full    string `json:"full"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Thumb   string `json:"thumb"`
	} `json:"urls"`
	User struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"user"`
	Links struct {
		Download string `json:"download"`
		HTML     string `json:"html"`
	} `json:"links"`
}

// TenorGIF represents a GIF from Tenor.
type TenorGIF struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Media []struct {
		GIF struct {
			URL string `json:"url"`
		} `json:"gif"`
		TinyGIF struct {
			URL string `json:"url"`
		} `json:"tinygif"`
	} `json:"media"`
}

// Search performs a search across Linkvite.
func (s *SearchService) Search(ctx context.Context, input *SearchInput) (*SearchResult, error) {
	v := url.Values{}
	v.Set("q", input.Query)
	if input.Type != "" {
		v.Set("type", input.Type)
	}
	if input.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", input.Page))
	}
	if input.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", input.Limit))
	}

	var result SearchResult
	if err := s.client.get(ctx, "/search?"+v.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchUnsplash searches for photos on Unsplash.
func (s *SearchService) SearchUnsplash(ctx context.Context, query string, page int64) ([]*UnsplashPhoto, error) {
	v := url.Values{}
	v.Set("q", query)
	if page > 0 {
		v.Set("page", fmt.Sprintf("%d", page))
	}

	var result []*UnsplashPhoto
	if err := s.client.get(ctx, "/search/unsplash?"+v.Encode(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TrackUnsplashDownload tracks a download for Unsplash attribution.
func (s *SearchService) TrackUnsplashDownload(ctx context.Context, photoID string) error {
	return s.client.post(ctx, "/search/unsplash/track", map[string]string{"photo_id": photoID}, nil)
}

// SearchTenor searches for GIFs on Tenor.
func (s *SearchService) SearchTenor(ctx context.Context, query string) ([]*TenorGIF, error) {
	v := url.Values{}
	v.Set("q", query)

	var result []*TenorGIF
	if err := s.client.get(ctx, "/search/tenor?"+v.Encode(), &result); err != nil {
		return nil, err
	}
	return result, nil
}
