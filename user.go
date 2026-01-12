package linkvite

import "context"

// UserService handles user operations.
type UserService struct {
	client *Client
}

// EditUserInput represents input for editing user details.
type EditUserInput struct {
	Name           *string `json:"name,omitempty"`
	Username       *string `json:"username,omitempty"`
	Avatar         *string `json:"avatar,omitempty"`
	FolderName     *string `json:"folder_name,omitempty"`
	PrivateAccount *bool   `json:"private_account,omitempty"`
}

// EditSettingsInput represents input for editing user settings.
type EditSettingsInput struct {
	PushEnabled  *bool `json:"push_enabled,omitempty"`
	EmailEnabled *bool `json:"email_enabled,omitempty"`
	InAppEnabled *bool `json:"in_app_enabled,omitempty"`
}

// Get retrieves the current authenticated user.
func (s *UserService) Get(ctx context.Context) (*User, error) {
	var result User
	if err := s.client.get(ctx, "/user", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates the current user's profile.
func (s *UserService) Edit(ctx context.Context, input *EditUserInput) (*User, error) {
	var result User
	if err := s.client.patch(ctx, "/user", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EditSettings updates the current user's notification settings.
func (s *UserService) EditSettings(ctx context.Context, input *EditSettingsInput) (*User, error) {
	var result User
	if err := s.client.patch(ctx, "/user/settings", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStorageUsage retrieves the user's storage usage.
func (s *UserService) GetStorageUsage(ctx context.Context) (*StorageUsage, error) {
	var result StorageUsage
	if err := s.client.get(ctx, "/user/storage", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDetailedStorageUsage retrieves detailed storage usage breakdown.
func (s *UserService) GetDetailedStorageUsage(ctx context.Context, opts *ListOptions) (*ListBookmarksWithStorageResult, error) {
	var result ListBookmarksWithStorageResult
	if err := s.client.get(ctx, "/user/storage/detailed"+opts.toQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EmptyTrash permanently deletes all trashed bookmarks.
func (s *UserService) EmptyTrash(ctx context.Context) error {
	return s.client.delete(ctx, "/user/trash")
}
