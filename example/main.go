// Example: list your bookmarks and save one.
//
//	LINKVITE_API_KEY=link_... go run ./example
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/linkvite/go/v2"
)

func main() {
	ctx := context.Background()
	client, err := linkvite.NewClient(os.Getenv("LINKVITE_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	me, err := client.User.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("signed in as", me.Username)

	b, err := client.Bookmarks.Create(ctx, &linkvite.CreateBookmarkRequest{URL: "https://go.dev/blog", Tags: "go"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("saved", b.ID)

	items := linkvite.Paginate(ctx, func(ctx context.Context, page int) ([]linkvite.BookmarkResponse, *linkvite.Pagination, error) {
		r, err := client.Bookmarks.GetAll(ctx, &linkvite.BookmarksGetAllParams{Page: linkvite.Ptr(page), Limit: linkvite.Ptr(100)})
		if err != nil {
			return nil, nil, err
		}
		return r.Bookmarks, r.Pagination, nil
	})
	for b, err := range items {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("-", b.Title, b.URL)
	}
}
