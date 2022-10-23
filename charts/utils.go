package charts

import (
	"context"
	"fmt"
	"time"

	"github.com/bisohns/saido/config"
)

// Periodic executes the provided closure periodically every interval.
// Exits when the context expires.
func Periodic(ctx context.Context, interval time.Duration, fn func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := fn(); err != nil {
				//FIXME: Best possible way to handle the
				// Update function failure (a retry mechanism)?
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Paginate helps to split hosts into multiple pages
func Paginate(hosts []config.Host, perPage int) [][]config.Host {
	if perPage < 1 {
		panic("A page must have at least 1")
	}

	paginated := [][]config.Host{}
	inner := []config.Host{}
	for ind, host := range hosts {
		ind = ind + 1
		inner = append(inner, host)
		// if a page is completed or we have reached
		// the end of hosts
		if ind%perPage == 0 || ind == len(hosts) {
			paginated = append(paginated, inner)
			inner = []config.Host{}
		}
	}
	return paginated
}

// Next go to next page without exceeding page length
func Next(current, pageLength int) int {
	if current+1 >= pageLength {
		return current
	}
	return current + 1
}

// Prev go to prev page without going below first page
func Prev(current, pageLength int) int {
	if current-1 <= 0 {
		return 0
	}
	return current - 1
}
