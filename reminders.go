package linkvite

import (
	"context"
	"time"
)

// RemindersService handles reminder operations.
type RemindersService struct {
	client *Client
}

// CreateReminderInput represents input for creating a reminder.
type CreateReminderInput struct {
	BookmarkID string    `json:"bookmark_id"`
	RemindAt   time.Time `json:"remind_at"`
	Note       string    `json:"note,omitempty"`
}

// EditReminderInput represents input for editing a reminder.
type EditReminderInput struct {
	RemindAt  *time.Time `json:"remind_at,omitempty"`
	Note      *string    `json:"note,omitempty"`
	Completed *bool      `json:"completed,omitempty"`
}

// BatchEditRemindersInput represents input for batch editing reminders.
type BatchEditRemindersInput struct {
	Reminders []string `json:"reminders"`
	Completed *bool    `json:"completed,omitempty"`
}

// BatchDeleteRemindersInput represents input for batch deleting reminders.
type BatchDeleteRemindersInput struct {
	Reminders []string `json:"reminders"`
}

// ListRemindersResult represents a paginated list of reminders.
type ListRemindersResult struct {
	Reminders  []*ReminderWithMeta `json:"reminders"`
	Pagination *Pagination         `json:"pagination"`
}

// List retrieves all reminders for the current user.
func (s *RemindersService) List(ctx context.Context, opts *ListOptions) (*ListRemindersResult, error) {
	var result ListRemindersResult
	if err := s.client.get(ctx, "/reminders"+opts.toQuery(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a single reminder by ID.
func (s *RemindersService) Get(ctx context.Context, id string) (*ReminderWithMeta, error) {
	var result ReminderWithMeta
	if err := s.client.get(ctx, "/reminders/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new reminder.
func (s *RemindersService) Create(ctx context.Context, input *CreateReminderInput) (*Reminder, error) {
	var result Reminder
	if err := s.client.post(ctx, "/reminders", input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates a reminder.
func (s *RemindersService) Edit(ctx context.Context, id string, input *EditReminderInput) (*Reminder, error) {
	var result Reminder
	if err := s.client.patch(ctx, "/reminders/"+id, input, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a reminder.
func (s *RemindersService) Delete(ctx context.Context, id string) error {
	return s.client.delete(ctx, "/reminders/"+id)
}

// BatchEdit updates multiple reminders at once.
func (s *RemindersService) BatchEdit(ctx context.Context, input *BatchEditRemindersInput) error {
	return s.client.post(ctx, "/reminders/batch-edit", input, nil)
}

// BatchDelete deletes multiple reminders at once.
func (s *RemindersService) BatchDelete(ctx context.Context, input *BatchDeleteRemindersInput) error {
	return s.client.post(ctx, "/reminders/batch-delete", input, nil)
}

// Clear deletes all reminders.
func (s *RemindersService) Clear(ctx context.Context) error {
	return s.client.post(ctx, "/reminders/clear", nil, nil)
}
