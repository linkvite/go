# Linkvite Go SDK

🦫 Official Go SDK for [Linkvite](https://linkvite.io) - a bookmarking service for saving, organizing, and sharing bookmarks.

## Installation

```bash
go get github.com/linkvite/go
```

## Quick Start

### With API Key

```go
package main

import (
    "context"
    "fmt"
    "log"

    linkvite "github.com/linkvite/go"
)

func main() {
    client, err := linkvite.NewClient("link_your-api-key")
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()

    // Get current user
    user, err := client.User.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Hello, %s!\n", user.Name)

    // Create a bookmark
    bookmark, err := client.Bookmarks.Create(ctx, &linkvite.CreateBookmarkInput{
        URL: "https://example.com",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created bookmark: %s\n", bookmark.Title)
}
```

### With Access Token

```go
// Using access token only
client, err := linkvite.NewClient("", linkvite.WithAccessToken("your-access-token"))

// Using access token with refresh token
client, err := linkvite.NewClientWithTokens("access-token", "refresh-token")

// Or using options
client, err := linkvite.NewClient("",
    linkvite.WithTokens("access-token", "refresh-token"),
)

// Refresh the access token when needed
err := client.RefreshAccessToken(ctx)
```

## Configuration

```go
client, err := linkvite.NewClient("link_your-api-key",
    linkvite.WithBaseURL("https://custom.api.url"),
    linkvite.WithHTTPClient(customHTTPClient),
    linkvite.WithUserAgent("my-app/1.0"),
)
```

## Services

### Bookmarks

```go
// List bookmarks
result, err := client.Bookmarks.List(ctx, &linkvite.ListOptions{
    Page:  1,
    Limit: 20,
})

// Get a single bookmark
bookmark, err := client.Bookmarks.Get(ctx, "bmk_abc123")

// Create a bookmark
bookmark, err := client.Bookmarks.Create(ctx, &linkvite.CreateBookmarkInput{
    URL:         "https://example.com",
    Title:       "Example Site",
    Description: "An example website",
    Collection:  "coll_xyz789",
    Tags:        "example,test",
})

// Edit a bookmark
bookmark, err := client.Bookmarks.Edit(ctx, "bmk_abc123", &linkvite.EditBookmarkInput{
    Name: "New Title",
})

// Delete a bookmark
err := client.Bookmarks.Delete(ctx, "bmk_abc123")

// Toggle star
err := client.Bookmarks.ToggleStar(ctx, "bmk_abc123")

// Check if bookmark exists
result, err := client.Bookmarks.Exists(ctx, &linkvite.BookmarkExistsInput{
    URL: "https://example.com",
})
```

### Collections

```go
// List collections
collections, err := client.Collections.List(ctx, nil)

// Get a collection
collection, err := client.Collections.Get(ctx, "coll_xyz789")

// Create a collection
collection, err := client.Collections.Create(ctx, &linkvite.CreateCollectionInput{
    Name:        "My Collection",
    Description: "A collection of bookmarks",
})

// Edit a collection
collection, err := client.Collections.Edit(ctx, "coll_xyz789", &linkvite.EditCollectionInput{
    Name:      "New Name",
    IsPrivate: linkvite.Ptr(false),
})

// Delete a collection
err := client.Collections.Delete(ctx, "coll_xyz789")

// Get collection members
members, err := client.Collections.GetMembers(ctx, "coll_xyz789")

// Invite a user to a collection
err := client.Collections.InviteUser(ctx, "coll_xyz789", &linkvite.InviteUserInput{
    Identifier: "user@example.com",
    Role:       linkvite.RoleLevelViewer,
})

// Leave a collection
err := client.Collections.Leave(ctx, "coll_xyz789")
```

### User

```go
// Get current user
user, err := client.User.Get(ctx)

// Edit user profile
user, err := client.User.Edit(ctx, &linkvite.EditUserInput{
    Name: linkvite.Ptr("New Name"),
})

// Edit user settings
user, err := client.User.EditSettings(ctx, &linkvite.EditSettingsInput{
    PushEnabled: linkvite.Ptr(true),
})

// Get storage usage
storage, err := client.User.GetStorageUsage(ctx)

// Empty trash
err := client.User.EmptyTrash(ctx)
```

### Comments

```go
// List comments on a bookmark
comments, err := client.Comments.List(ctx, "bmk_abc123")

// Create a comment
comment, err := client.Comments.Create(ctx, &linkvite.CreateCommentInput{
    BookmarkID: "bmk_abc123",
    Content:    "Great article!",
})

// Delete a comment
err := client.Comments.Delete(ctx, "cmt_xyz789")

// Toggle like on a comment
err := client.Comments.ToggleLike(ctx, "cmt_xyz789")

// Get likes on a comment
likes, err := client.Comments.GetLikes(ctx, "cmt_xyz789")

// Report a comment
err := client.Comments.Report(ctx, "cmt_xyz789")
```

### Highlights

```go
// List highlights on a bookmark
result, err := client.Highlights.List(ctx, "bmk_abc123")

// Create a highlight
result, err := client.Highlights.Create(ctx, &linkvite.CreateHighlightInput{
    ID:         "hil_xyz789",  // Client-generated ID
    BookmarkID: "bmk_abc123",
    Color:      linkvite.HighlightColor1,
    Start:      linkvite.Ptr(0),
    End:        linkvite.Ptr(50),
    Exact:      "The exact highlighted text",
    Note:       "My note about this",
})

// Available highlight colors
linkvite.HighlightColor1  // "lv-highlighter-1"
linkvite.HighlightColor2  // "lv-highlighter-2"
linkvite.HighlightColor3  // "lv-highlighter-3"
linkvite.HighlightColor4  // "lv-highlighter-4"
linkvite.HighlightColor5  // "lv-highlighter-5"

// Edit a highlight
highlight, err := client.Highlights.Edit(ctx, "hil_xyz789", &linkvite.EditHighlightInput{
    Note:  "Updated note",
    Color: linkvite.HighlightColor2,
})

// Delete a highlight
err := client.Highlights.Delete(ctx, "hil_xyz789")
```

### Reminders

```go
// List reminders
result, err := client.Reminders.List(ctx, &linkvite.ListOptions{})

// Get a reminder
reminder, err := client.Reminders.Get(ctx, "rem_abc123")

// Create a reminder
reminder, err := client.Reminders.Create(ctx, &linkvite.CreateReminderInput{
    BookmarkID: "bmk_abc123",
    RemindAt:   time.Now().Add(24 * time.Hour),
    Note:       "Read this article",
})

// Edit a reminder
reminder, err := client.Reminders.Edit(ctx, "rem_abc123", &linkvite.EditReminderInput{
    Completed: linkvite.Ptr(true),
})

// Delete a reminder
err := client.Reminders.Delete(ctx, "rem_abc123")

// Batch operations
err := client.Reminders.BatchEdit(ctx, &linkvite.BatchEditRemindersInput{
    Reminders: []string{"rem_1", "rem_2"},
    Completed: linkvite.Ptr(true),
})

err := client.Reminders.BatchDelete(ctx, &linkvite.BatchDeleteRemindersInput{
    Reminders: []string{"rem_1", "rem_2"},
})

err := client.Reminders.Clear(ctx)
```

### Invites

```go
// List pending invites
invites, err := client.Invites.List(ctx)

// Get an invite
invite, err := client.Invites.Get(ctx, "inv_abc123")

// Accept an invite
err := client.Invites.Accept(ctx, "inv_abc123")

// Decline an invite
err := client.Invites.Decline(ctx, "inv_abc123")

// Batch operations
err := client.Invites.BatchAccept(ctx, &linkvite.BatchInviteInput{
    Invites: []string{"inv_1", "inv_2"},
})

err := client.Invites.BatchDecline(ctx, &linkvite.BatchInviteInput{
    Invites: []string{"inv_1", "inv_2"},
})

// Decline all invites
err := client.Invites.DeclineAll(ctx)
```

### Search

```go
// Search Linkvite
results, err := client.Search.Search(ctx, &linkvite.SearchInput{
    Query: "golang",
    Type:  "bookmarks", // "bookmarks", "collections", or "users"
    Page:  1,
    Limit: 20,
})

// Search Unsplash for images
photos, err := client.Search.SearchUnsplash(ctx, "nature", 1)

// Track Unsplash download (for attribution)
err := client.Search.TrackUnsplashDownload(ctx, "photo_id")

// Search Tenor for GIFs
gifs, err := client.Search.SearchTenor(ctx, "coding")
```

### RSS Feeds

```go
// List RSS feeds
feeds, err := client.RSSFeeds.List(ctx, &linkvite.ListOptions{})

// Create an RSS feed subscription
feed, err := client.RSSFeeds.Create(ctx, &linkvite.CreateRSSFeedInput{
    URL:          "https://example.com/feed.xml",
    CollectionID: "coll_xyz789",
})

// Get RSS feed items
items, err := client.RSSFeeds.GetItems(ctx, "rss_abc123")

// Edit an RSS feed
feed, err := client.RSSFeeds.Edit(ctx, "rss_abc123", &linkvite.EditRSSFeedInput{
    Title: "New Title",
})

// Batch delete RSS feeds
err := client.RSSFeeds.BatchDelete(ctx, &linkvite.BatchDeleteRSSFeedsInput{
    Feeds: []string{"rss_1", "rss_2"},
})
```

### API Keys

```go
// List API keys
keys, err := client.APIKeys.List(ctx)

// Create an API key
result, err := client.APIKeys.Create(ctx, &linkvite.CreateAPIKeyInput{
    Name:   "My API Key",
    Scopes: []string{"read", "write"},
})
// result.Key contains the full API key (only available once)

// Edit an API key
key, err := client.APIKeys.Edit(ctx, "key_abc123", &linkvite.EditAPIKeyInput{
    Name: "Updated Name",
})
```

### Parser

```go
// Parse a URL to get its metadata
parsed, err := client.Parser.ParseLink(ctx, "https://example.com")

// Import parsed items as bookmarks
err := client.Parser.Import(ctx, &linkvite.ImportParsedItemsInput{
    Items:      []linkvite.ParsedLink{...},
    Collection: "coll_xyz789",
})
```

## Error Handling

```go
bookmark, err := client.Bookmarks.Get(ctx, "invalid-id")
if err != nil {
    if apiErr, ok := err.(*linkvite.Error); ok {
        if apiErr.IsNotFound() {
            fmt.Println("Bookmark not found")
        } else if apiErr.IsUnauthorized() {
            fmt.Println("Invalid authentication")
        } else if apiErr.IsForbidden() {
            fmt.Println("Access forbidden")
        } else if apiErr.IsRateLimited() {
            fmt.Println("Rate limited, try again later")
        } else {
            fmt.Printf("API error: %s\n", apiErr.Message)
        }
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Helper Functions

```go
// Use Ptr() to create pointers for optional fields
input := &linkvite.EditBookmarkInput{
    Description: linkvite.Ptr("New description"),
    Starred:     linkvite.Ptr(true),
}
```

## License

MIT License - see [LICENSE](LICENSE) for details.
