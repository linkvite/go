package linkvite

import "context"

// InvitesService handles invite operations.
type InvitesService struct {
	client *Client
}

// BatchInviteInput represents input for batch invite operations.
type BatchInviteInput struct {
	Invites []string `json:"invites"`
}

// List retrieves all pending invites for the current user.
func (s *InvitesService) List(ctx context.Context) ([]*InviteWithMeta, error) {
	var result []*InviteWithMeta
	if err := s.client.get(ctx, "/invites", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Get retrieves a single invite by ID.
func (s *InvitesService) Get(ctx context.Context, id string) (*InviteWithMeta, error) {
	var result InviteWithMeta
	if err := s.client.get(ctx, "/invites/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Accept accepts an invite.
func (s *InvitesService) Accept(ctx context.Context, id string) error {
	return s.client.post(ctx, "/invites/"+id+"/accept", nil, nil)
}

// Decline declines an invite.
func (s *InvitesService) Decline(ctx context.Context, id string) error {
	return s.client.post(ctx, "/invites/"+id+"/decline", nil, nil)
}

// BatchAccept accepts multiple invites at once.
func (s *InvitesService) BatchAccept(ctx context.Context, input *BatchInviteInput) error {
	return s.client.post(ctx, "/invites/batch-accept", input, nil)
}

// BatchDecline declines multiple invites at once.
func (s *InvitesService) BatchDecline(ctx context.Context, input *BatchInviteInput) error {
	return s.client.post(ctx, "/invites/batch-decline", input, nil)
}

// DeclineAll declines all pending invites.
func (s *InvitesService) DeclineAll(ctx context.Context) error {
	return s.client.post(ctx, "/invites/decline/all", nil, nil)
}
