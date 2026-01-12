package linkvite

import "context"

// CommentsService handles comment operations.
type CommentsService struct {
	client *Client
}

// CreateCommentInput represents input for creating a comment.
type CreateCommentInput struct {
	BookmarkID string `json:"bookmark_id"`
	Content    string `json:"content"`
}

// List retrieves all comments for a bookmark.
func (s *CommentsService) List(ctx context.Context, bookmarkID string) ([]*CommentWithMeta, error) {
	var result []*CommentWithMeta
	if err := s.client.get(ctx, "/comments/"+bookmarkID, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create creates a new comment on a bookmark.
func (s *CommentsService) Create(ctx context.Context, input *CreateCommentInput) (*Comment, error) {
	var result Comment
	if err := s.client.post(ctx, "/comments", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a comment.
func (s *CommentsService) Delete(ctx context.Context, id string) error {
	return s.client.delete(ctx, "/comments/"+id)
}

// ToggleLike toggles the like status of a comment.
func (s *CommentsService) ToggleLike(ctx context.Context, id string) error {
	return s.client.post(ctx, "/comments/"+id+"/like", nil, nil)
}

// GetLikes retrieves users who liked a comment.
func (s *CommentsService) GetLikes(ctx context.Context, id string) ([]*PublicProfile, error) {
	var result []*PublicProfile
	if err := s.client.get(ctx, "/comments/"+id+"/likes", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Report reports a comment for moderation.
func (s *CommentsService) Report(ctx context.Context, id string) error {
	return s.client.post(ctx, "/comments/"+id+"/report", nil, nil)
}
