package linkvite

import (
	"fmt"
	"net/url"
	"time"
)

// UserStatus represents the status of a user account.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// AccountType represents the subscription type of a user.
type AccountType string

const (
	AccountTypeFree       AccountType = "free"
	AccountTypePro        AccountType = "pro"
	AccountTypeEnterprise AccountType = "enterprise"
)

// User represents a Linkvite user.
type User struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	Username       string      `json:"username"`
	Avatar         string      `json:"avatar"`
	EmailVerified  bool        `json:"email_verified"`
	AccountType    AccountType `json:"account_type"`
	Status         UserStatus  `json:"status"`
	Verified       bool        `json:"verified"`
	FolderName     string      `json:"folder_name"`
	PrivateAccount bool        `json:"private_account"`
	LastLoginAt    *time.Time  `json:"last_login_at"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	PushEnabled    bool        `json:"push_enabled"`
	EmailEnabled   bool        `json:"email_enabled"`
	InAppEnabled   bool        `json:"in_app_enabled"`
}

// PublicProfile represents a public user profile.
type PublicProfile struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Username       string      `json:"username"`
	Avatar         string      `json:"avatar"`
	AccountType    AccountType `json:"account_type"`
	Status         UserStatus  `json:"status"`
	Verified       bool        `json:"verified"`
	CreatedAt      time.Time   `json:"created_at"`
	PrivateAccount bool        `json:"private_account"`
	Role           RoleLevel   `json:"role,omitempty"`
}

type BookmarkMimeType int

const (
	BookmarkMimeTypeHTML BookmarkMimeType = iota
	BookmarkMimeTypePlain
	BookmarkMimeTypePDF
	BookmarkMimeTypeZIP
	BookmarkMimeTypeDOCX
	BookmarkMimeTypePPTX
	BookmarkMimeTypeXLSX
	BookmarkMimeTypeJPG
	BookmarkMimeTypeJPEG
	BookmarkMimeTypePNG
	BookmarkMimeTypeGIF
	BookmarkMimeTypeWEBP
	BookmarkMimeTypeSVG
	BookmarkMimeTypeBMP
	BookmarkMimeTypeMP3
	BookmarkMimeTypeWAV
	BookmarkMimeTypeOGG
	BookmarkMimeTypeM4A
	BookmarkMimeTypeMP4Audio
	BookmarkMimeTypeMP4Video
	BookmarkMimeTypeWEBM
)

// BookmarkStatus represents the status of a bookmark.
type BookmarkStatus string

const (
	BookmarkStatusActive  BookmarkStatus = "active"
	BookmarkStatusTrashed BookmarkStatus = "trashed"
)

// BookmarkArchiveStatus represents the archive status of a bookmark.
type BookmarkArchiveStatus string

const (
	BookmarkArchiveStatusPending BookmarkArchiveStatus = "pending"
	BookmarkArchiveStatusSuccess BookmarkArchiveStatus = "success"
	BookmarkArchiveStatusFailed  BookmarkArchiveStatus = "failed"
)

// Bookmark represents a Linkvite bookmark.
type Bookmark struct {
	ID            string                `json:"id"`
	UserID        string                `json:"user_id"`
	Title         string                `json:"title"`
	Description   string                `json:"description,omitempty"`
	CollectionID  *string               `json:"collection_id"`
	Status        BookmarkStatus        `json:"status"`
	IsBroken      bool                  `json:"is_broken"`
	AllowComments bool                  `json:"allow_comments"`
	Tags          string                `json:"tags"`
	Icon          string                `json:"icon"`
	Thumbnail     string                `json:"thumbnail"`
	Processing    bool                  `json:"processing"`
	Type          string                `json:"type"`
	URL           string                `json:"url"`
	Size          float64               `json:"size"`
	Views         int                   `json:"views"`
	UpdatedByID   *string               `json:"updated_by_id"`
	UpdatedAt     time.Time             `json:"updated_at"`
	CreatedAt     time.Time             `json:"created_at"`
	LastOpenedAt  *time.Time            `json:"last_opened_at"`
	ArchiveStatus BookmarkArchiveStatus `json:"archive_status"`
	ArchiveSize   float64               `json:"archive_size"`
	RSSFeedID     *string               `json:"rss_feed_id,omitempty"`
	LocalID       *string               `json:"local_id"`
	MimeTypeID    BookmarkMimeType      `json:"mime_type_id"`
}

// BookmarkWithMeta represents a bookmark with additional metadata.
type BookmarkWithMeta struct {
	Bookmark

	OwnerID        string    `json:"owner_id"`
	OwnerName      string    `json:"owner_name"`
	OwnerUsername  string    `json:"owner_username"`
	OwnerAvatar    string    `json:"owner_avatar"`
	CollectionName *string   `json:"collection_name,omitempty"`
	Starred        bool      `json:"starred"`
	Role           RoleLevel `json:"role"`
}

// BookmarkStorage represents a bookmark with storage information.
type BookmarkStorage struct {
	ID             string           `json:"id"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Thumbnail      string           `json:"thumbnail"`
	Icon           string           `json:"icon"`
	URL            string           `json:"url"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	UserID         string           `json:"user_id"`
	CollectionID   *string          `json:"collection_id"`
	CollectionName *string          `json:"collection_name"`
	Size           float64          `json:"size"`
	MimeTypeID     BookmarkMimeType `json:"mime_type_id"`
	Status         BookmarkStatus   `json:"status"`
}

// BookmarkReader represents readable content from a bookmark.
type BookmarkReader struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	ByLine      string `json:"byline"`
	Content     string `json:"content"`
	TextContent string `json:"text_content"`
}

// RoleLevel represents the role level in a collection.
type RoleLevel string

const (
	RoleLevelAdmin  RoleLevel = "admin"
	RoleLevelViewer RoleLevel = "viewer"
)

// CollectionStatus represents the status of a collection.
type CollectionStatus string

const (
	CollectionStatusActive   CollectionStatus = "active"
	CollectionStatusInactive CollectionStatus = "inactive"
)

// Collection represents a Linkvite collection.
type Collection struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Name            string           `json:"name"`
	Description     string           `json:"description,omitempty"`
	Status          CollectionStatus `json:"status"`
	IsPrivate       bool             `json:"is_private"`
	AllowPublicJoin bool             `json:"allow_public_join"`
	Thumbnail       string           `json:"thumbnail"`
	Icon            string           `json:"icon"`
	Views           int              `json:"views"`
	ParentID        *string          `json:"parent_id"`
	UpdatedByID     *string          `json:"updated_by_id"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	LastOpenedAt    *time.Time       `json:"last_opened_at"`
	Role            RoleLevel        `json:"role"`
}

// CollectionWithMeta represents a collection with additional metadata.
type CollectionWithMeta struct {
	Collection

	OwnerID        string     `json:"owner_id"`
	OwnerName      string     `json:"owner_name"`
	OwnerUsername  string     `json:"owner_username"`
	OwnerAvatar    string     `json:"owner_avatar"`
	BookmarksCount int        `json:"bookmarks_count"`
	Role           *RoleLevel `json:"role"`
}

// CollectionMember represents a member of a collection.
type CollectionMember struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Role     RoleLevel `json:"role"`
}

// Comment represents a comment on a bookmark.
type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentWithMeta represents a comment with user information.
type CommentWithMeta struct {
	Comment

	UserName     string `json:"user_name"`
	UserUsername string `json:"user_username"`
	UserAvatar   string `json:"user_avatar"`
	LikesCount   int    `json:"likes_count"`
	Liked        bool   `json:"liked"`
}

// Highlight represents a text highlight in a bookmark.
type Highlight struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	BookmarkID string    `json:"bookmark_id"`
	Note       string    `json:"note"`
	Exact      string    `json:"exact"`
	Start      int       `json:"start"`
	End        int       `json:"end"`
	Color      string    `json:"color"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// HighlightedBody represents the body of a highlighted annotation.
type HighlightedBody struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Purpose string `json:"purpose"`
}

// HighlightedTargetSelector represents a selector in a highlight target.
type HighlightedTargetSelector struct {
	Type  string `json:"type"`
	Exact string `json:"exact,omitempty"`
	Start int    `json:"start,omitempty"`
	End   int    `json:"end,omitempty"`
}

// HighlightedTarget represents the target of a highlighted annotation.
type HighlightedTarget struct {
	Selector []HighlightedTargetSelector `json:"selector"`
}

// Highlighted represents a highlight in annotation format.
type Highlighted struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	BookmarkID string            `json:"bookmark_id"`
	Note       string            `json:"note"`
	Type       string            `json:"type"`
	Body       []HighlightedBody `json:"body"`
	Target     HighlightedTarget `json:"target"`
}

// Reminder represents a bookmark reminder.
type Reminder struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	BookmarkID string    `json:"bookmark_id"`
	RemindAt   time.Time `json:"remind_at"`
	Note       string    `json:"note,omitempty"`
	Completed  bool      `json:"completed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ReminderWithMeta represents a reminder with bookmark information.
type ReminderWithMeta struct {
	Reminder

	BookmarkTitle     string `json:"bookmark_title"`
	BookmarkURL       string `json:"bookmark_url"`
	BookmarkThumbnail string `json:"bookmark_thumbnail"`
}

// Invite represents a collection invite.
type Invite struct {
	ID           string    `json:"id"`
	CollectionID string    `json:"collection_id"`
	InviterID    string    `json:"inviter_id"`
	InviteeID    string    `json:"invitee_id"`
	Role         RoleLevel `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// InviteWithMeta represents an invite with additional information.
type InviteWithMeta struct {
	Invite

	CollectionName      string `json:"collection_name"`
	CollectionThumbnail string `json:"collection_thumbnail"`
	InviterName         string `json:"inviter_name"`
	InviterUsername     string `json:"inviter_username"`
	InviterAvatar       string `json:"inviter_avatar"`
}

// RSSFeed represents an RSS feed subscription.
type RSSFeed struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	Description  string     `json:"description,omitempty"`
	CollectionID *string    `json:"collection_id,omitempty"`
	LastFetchAt  *time.Time `json:"last_fetch_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// RSSFeedWithMeta represents an RSS feed with additional information.
type RSSFeedWithMeta struct {
	RSSFeed

	CollectionName string `json:"collection_name,omitempty"`
	ItemCount      int    `json:"item_count"`
}

// RSSFeedItem represents an item from an RSS feed.
type RSSFeedItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Published   string `json:"published"`
	Author      string `json:"author"`
}

// APIKey represents an API key.
type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Scopes     []string   `json:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

// StorageUsage represents user storage usage.
type StorageUsage struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

// DetailedStorageUsage represents detailed storage usage.
type DetailedStorageUsage struct {
	Bookmarks float64 `json:"bookmarks"`
	Archives  float64 `json:"archives"`
	Files     float64 `json:"files"`
	Total     float64 `json:"total"`
	Limit     float64 `json:"limit"`
}

// ParsedLink represents metadata parsed from a URL.
type ParsedLink struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Favicon     string `json:"favicon"`
	Type        string `json:"type"`
	SiteName    string `json:"site_name"`
}

// Pagination represents pagination information.
type Pagination struct {
	Query    string `json:"q,omitempty"`
	Total    int64  `json:"total"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
	LastPage int64  `json:"last_page"`
}

// ListOptions represents common options for list operations.
type ListOptions struct {
	Page   int64
	Limit  int64
	Query  string
	Status string
	Sort   string
}

func (o *ListOptions) toQuery() string {
	if o == nil {
		return ""
	}
	v := url.Values{}
	if o.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", o.Page))
	}
	if o.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", o.Limit))
	}
	if o.Query != "" {
		v.Set("q", o.Query)
	}
	if o.Status != "" {
		v.Set("status", o.Status)
	}
	if o.Sort != "" {
		v.Set("sort", o.Sort)
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}
