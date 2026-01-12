package linkvite

import (
	"context"
	"net/url"
)

// ParserService handles link parsing operations.
type ParserService struct {
	client *Client
}

// ParseLinkInput represents input for parsing a link.
type ParseLinkInput struct {
	URL string `json:"url"`
}

// ImportParsedItemsInput represents input for importing parsed items.
type ImportParsedItemsInput struct {
	Items      []ParsedLink `json:"items"`
	Collection string       `json:"collection,omitempty"`
}

// ParseLink parses a URL and returns its metadata.
func (s *ParserService) ParseLink(ctx context.Context, linkURL string) (*ParsedLink, error) {
	v := url.Values{}
	v.Set("url", linkURL)

	var result ParsedLink
	if err := s.client.get(ctx, "/parser?"+v.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Import imports parsed items as bookmarks.
func (s *ParserService) Import(ctx context.Context, input *ImportParsedItemsInput) error {
	return s.client.post(ctx, "/parser/import", input, nil)
}
