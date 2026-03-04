package infrastructure

import (
	"github.com/redis/go-redis/v9"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/sqldb"
	"gorm.io/gorm"
)

// CreateRedisCon initializes and returns a Redis client.
func CreateRedisCon() *redis.Client {
	redisClient, err := redisPkg.NewClient("")
	common.HandlerError(err)

	return redisClient
}

// CreateSQLDB initializes and returns a GORM DB client without performing migrations.
// It returns the initialized GORM DB instance.
// Returns:
//   - *gorm.DB: A pointer to the initialized GORM DB instance
func CreateSQLDB() *gorm.DB {
	dbClient, err := sqldb.NewClient("")
	common.HandlerError(err)

	return dbClient
}
