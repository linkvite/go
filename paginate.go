package linkvite

import (
	"context"
	"iter"
)

// Paginate yields every item of a paginated list, fetching pages on demand.
// page fetches one page and returns its items and pagination.
//
//	items := linkvite.Paginate(ctx, func(ctx context.Context, page int) ([]BookmarkResponse, *Pagination, error) {
//		r, err := client.Bookmarks.GetAll(ctx, &linkvite.BookmarksGetAllParams{Page: linkvite.Ptr(page), Limit: linkvite.Ptr(100)})
//		if err != nil {
//			return nil, nil, err
//		}
//		return r.Bookmarks, r.Pagination, nil
//	})
//	for b, err := range items {
//		if err != nil { ... }
//		fmt.Println(b.Title)
//	}
func Paginate[T any](ctx context.Context, page func(ctx context.Context, page int) ([]T, *Pagination, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for n := 1; ; n++ {
			items, p, err := page(ctx, n)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, it := range items {
				if !yield(it, nil) {
					return
				}
			}
			if len(items) == 0 || p == nil || int64(n) >= p.LastPage {
				return
			}
		}
	}
}
