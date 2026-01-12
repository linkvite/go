package linkvite

import "context"

// APIKeysService handles API key operations.
type APIKeysService struct {
	client *Client
}

// CreateAPIKeyInput represents input for creating an API key.
type CreateAPIKeyInput struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes,omitempty"`
}

// EditAPIKeyInput represents input for editing an API key.
type EditAPIKeyInput struct {
	Name   string   `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

// CreateAPIKeyResult represents the result when creating an API key.
// The Key field contains the full API key and is only returned once.
type CreateAPIKeyResult struct {
	APIKey

	Key string `json:"key"`
}

// List retrieves all API keys for the current user.
func (s *APIKeysService) List(ctx context.Context) ([]*APIKey, error) {
	var result []*APIKey
	if err := s.client.get(ctx, "/api-keys", &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Create creates a new API key.
// The returned CreateAPIKeyResult contains the full API key in the Key field.
// This is the only time the full key is available.
func (s *APIKeysService) Create(ctx context.Context, input *CreateAPIKeyInput) (*CreateAPIKeyResult, error) {
	var result CreateAPIKeyResult
	if err := s.client.post(ctx, "/api-keys", input, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Edit updates an API key.
func (s *APIKeysService) Edit(ctx context.Context, id string, input *EditAPIKeyInput) (*APIKey, error) {
	var result APIKey
	if err := s.client.patch(ctx, "/api-keys/"+id, input, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
