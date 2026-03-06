package cache

import (
	"context"
)

// DB defines the interface for cache database operations.
// It includes methods for setting, getting, and deleting cache data.
//
//go:generate mockery --name=Repository --output=./mocks --filename=cache.go --outpkg=mock_cache
type Repository interface {
	// DelCacheData deletes the cache data for a given group key.
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - groupKey: The group key under which the cache data is stored.
	//
	// Returns:
	//   - error: An error if the operation fails, otherwise nil.
	DelCacheData(ctx context.Context, groupKey string) error
}
