package linkvite

import "context"

// RSSFeedsService handles RSS feed operations.
type RSSFeedsService struct {
	client *Client
}

// CreateRSSFeedInput represents input for creating an RSS feed.
type CreateRSSFeedInput struct {
	URL          string `json:"url"`
	CollectionID string `json:"collection_id,omitempty"`
}

// EditRSSFeedInput represents input for editing an RSS feed.
type EditRSSFeedInput struct {
	Title        string  `json:"title,omitempty"`
	CollectionID *string `json:"collection_id,omitempty"`
}

// BatchDeleteRSSFeedsInput represents input for batch deleting RSS feeds.
type BatchDeleteRSSFeedsInput struct {
	Feeds []string `json:"feeds"`
}

// List retrieves all RSS feeds for the current user.
func (s *RSSFeedsService) List(ctx context.Context, opts *ListOptions) ([]*RSSFeedWithMeta, error) {
	var result []*RSSFeedWithMeta
	if err := s.client.get(ctx, "/rss-feeds"+opts.toQuery(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetItems retrieves items from an RSS feed.
func (s *RSSFeedsService) GetItems(ctx context.Context, id string) ([]*RSSFeedItem, error) {
	var result []*RSSFeedItem
	if err := s.client.get(ctx, "/rss-feeds/"+id+"/items", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create creates a new RSS feed subscription.
func (s *RSSFeedsService) Create(ctx context.Context, input *CreateRSSFeedInput) (*RSSFeed, error) {
	var result RSSFeed
	if err := s.client.post(ctx, "/rss-feeds", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates an RSS feed.
func (s *RSSFeedsService) Edit(ctx context.Context, id string, input *EditRSSFeedInput) (*RSSFeed, error) {
	var result RSSFeed
	if err := s.client.patch(ctx, "/rss-feeds/"+id, input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchDelete deletes multiple RSS feeds at once.
func (s *RSSFeedsService) BatchDelete(ctx context.Context, input *BatchDeleteRSSFeedsInput) error {
	return s.client.post(ctx, "/rss-feeds/batch-delete", input, nil)
}
