package concurrency

import (
	"context"
	"time"

	"github.com/shomali11/parallelizer"
)

// Client manage all concurrency function
type Client struct{}

// Parallelize parallelizes the function calls
func (c *Client) Parallelize(functions ...func()) error {
	return ParallelizeTimeout(0, functions...)
}

// ParallelizeTimeout parallelizes the function calls with a timeout
func (c *Client) ParallelizeTimeout(timeout time.Duration, functions ...func()) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return ParallelizeContext(ctx, functions...)
}

// ParallelizeContext parallelizes the function calls with a context
func (c *Client) ParallelizeContext(ctx context.Context, functions ...func()) error {
	group := parallelizer.NewGroup()
	for _, function := range functions {
		group.Add(function)
	}

	return group.Wait(parallelizer.WithContext(ctx))
}
