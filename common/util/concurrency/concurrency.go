package concurrency

import (
	"context"
	"time"
)

// Storage interface for concurrency package 
type Storage interface {
	Parallelize(functions ...func()) error
	ParallelizeTimeout(timeout time.Duration, functions ...func()) error
	ParallelizeContext(ctx context.Context, functions ...func()) error
}
