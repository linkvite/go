package linkvite

import "context"

// HighlightsService handles highlight operations.
type HighlightsService struct {
	client *Client
}

// HighlightColor represents the color of a highlight.
type HighlightColor string

const (
	// HighlightColor1 looks like a pink highlighter.
	HighlightColor1 HighlightColor = "lv-highlighter-1"

	// HighlightColor2 looks like a blue highlighter.
	HighlightColor2 HighlightColor = "lv-highlighter-2"

	// HighlightColor3 looks like a green highlighter.
	HighlightColor3 HighlightColor = "lv-highlighter-3"

	// HighlightColor4 looks like a purple highlighter.
	HighlightColor4 HighlightColor = "lv-highlighter-4"

	// HighlightColor5 looks like an orange highlighter.
	HighlightColor5 HighlightColor = "lv-highlighter-5"
)

// CreateHighlightInput represents input for creating a highlight.
type CreateHighlightInput struct {
	ID         string         `json:"id"`
	BookmarkID string         `json:"bookmark_id"`
	Color      HighlightColor `json:"color"`
	Start      *int           `json:"start"`
	End        *int           `json:"end"`
	Exact      string         `json:"exact"`
	Note       string         `json:"note,omitempty"`
}

// EditHighlightInput represents input for editing a highlight.
type EditHighlightInput struct {
	Note  string         `json:"note,omitempty"`
	Color HighlightColor `json:"color,omitempty"`
}

// CreateHighlightResult represents the result of creating a highlight.
type CreateHighlightResult struct {
	Highlight  *Highlighted `json:"highlight"`
	PreviousID string       `json:"previous_id"`
}

// ListHighlightsResult represents the result of listing highlights.
type ListHighlightsResult struct {
	Highlights []*Highlight `json:"highlights"`
	Bookmark   string       `json:"bookmark"`
}

// List retrieves all highlights for a bookmark.
func (s *HighlightsService) List(ctx context.Context, bookmarkID string) (*ListHighlightsResult, error) {
	var result ListHighlightsResult
	if err := s.client.get(ctx, "/highlights/"+bookmarkID, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new highlight.
func (s *HighlightsService) Create(ctx context.Context, input *CreateHighlightInput) (*CreateHighlightResult, error) {
	var result CreateHighlightResult
	if err := s.client.post(ctx, "/highlights", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates a highlight.
func (s *HighlightsService) Edit(ctx context.Context, id string, input *EditHighlightInput) (*Highlight, error) {
	var result Highlight
	if err := s.client.patch(ctx, "/highlights/"+id, input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a highlight.
func (s *HighlightsService) Delete(ctx context.Context, id string) error {
	return s.client.delete(ctx, "/highlights/"+id)
}
