package linkvite

import "context"

// CollectionsService handles collection operations.
type CollectionsService struct {
	client *Client
}

// CreateCollectionInput represents input for creating a collection.
type CreateCollectionInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parent      string `json:"parent,omitempty"`
}

// EditCollectionInput represents input for editing a collection.
type EditCollectionInput struct {
	Name            string  `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	IsPrivate       *bool   `json:"is_private,omitempty"`
	AllowPublicJoin *bool   `json:"allow_public_join,omitempty"`
	Icon            string  `json:"icon,omitempty"`
	Cover           string  `json:"cover,omitempty"`
}

// InviteUserInput represents input for inviting a user to a collection.
type InviteUserInput struct {
	Identifier string    `json:"identifier"`
	Role       RoleLevel `json:"role"`
}

// EditMemberInput represents input for editing a collection member.
type EditMemberInput struct {
	Role   RoleLevel `json:"role,omitempty"`
	Remove *bool     `json:"remove,omitempty"`
	UserID string    `json:"user_id"`
}

// MoveCollectionInput represents input for moving a collection.
type MoveCollectionInput struct {
	Parent *string `json:"parent,omitempty"`
	Remove *bool   `json:"remove,omitempty"`
}

// CollectionMembersResult represents the result for collection members.
type CollectionMembersResult struct {
	Role    RoleLevel           `json:"role"`
	Members []*CollectionMember `json:"members"`
	OwnerID string              `json:"owner_id"`
}

// CollectionPresenceResult represents collection presence info.
type CollectionPresenceResult struct {
	TotalMembers  int64 `json:"total_members"`
	OnlineMembers int   `json:"online_members"`
	OnlineClients int   `json:"online_clients"`
}

// List retrieves all collections for the current user.
func (s *CollectionsService) List(ctx context.Context, opts *ListOptions) ([]*CollectionWithMeta, error) {
	var result []*CollectionWithMeta
	if err := s.client.get(ctx, "/collections"+opts.toQuery(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Get retrieves a single collection by ID.
func (s *CollectionsService) Get(ctx context.Context, id string) (*CollectionWithMeta, error) {
	var result CollectionWithMeta
	if err := s.client.get(ctx, "/collections/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new collection.
func (s *CollectionsService) Create(ctx context.Context, input *CreateCollectionInput) (*Collection, error) {
	var result Collection
	if err := s.client.post(ctx, "/collections", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates a collection.
func (s *CollectionsService) Edit(ctx context.Context, id string, input *EditCollectionInput) (*Collection, error) {
	var result Collection
	if err := s.client.patch(ctx, "/collections/"+id, input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a collection.
func (s *CollectionsService) Delete(ctx context.Context, id string) error {
	return s.client.delete(ctx, "/collections/"+id)
}

// ToggleStar toggles the starred status of a collection.
func (s *CollectionsService) ToggleStar(ctx context.Context, id string) error {
	return s.client.post(ctx, "/collections/"+id+"/star", nil, nil)
}

// GetMembers retrieves the members of a collection.
func (s *CollectionsService) GetMembers(ctx context.Context, id string) (*CollectionMembersResult, error) {
	var result CollectionMembersResult
	if err := s.client.get(ctx, "/collections/"+id+"/members", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InviteUser invites a user to a collection.
func (s *CollectionsService) InviteUser(ctx context.Context, id string, input *InviteUserInput) error {
	return s.client.post(ctx, "/collections/"+id+"/send-invite", input, nil)
}

// EditMember updates a member's role or removes them.
func (s *CollectionsService) EditMember(ctx context.Context, id string, input *EditMemberInput) error {
	return s.client.patch(ctx, "/collections/"+id+"/edit-role", input, nil)
}

// Leave removes the current user from a collection.
func (s *CollectionsService) Leave(ctx context.Context, id string) error {
	return s.client.post(ctx, "/collections/"+id+"/leave", nil, nil)
}

// Move moves a collection to a new parent.
func (s *CollectionsService) Move(ctx context.Context, id string, input *MoveCollectionInput) error {
	return s.client.post(ctx, "/collections/"+id+"/move", input, nil)
}

// GetPresence retrieves online presence for a collection.
func (s *CollectionsService) GetPresence(ctx context.Context, id string) (*CollectionPresenceResult, error) {
	var result CollectionPresenceResult
	if err := s.client.get(ctx, "/collections/"+id+"/presence", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
