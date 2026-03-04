package data

import (
	"time"

	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	"gorm.io/gorm"
)

type BookmarkCommonTestDB struct {
	base
}

// Migrate migrates the database schema for the BookmarkCommonTestDB fixture.
//
// Returns:
//   - error: An error if migration fails, otherwise nil
func (b *BookmarkCommonTestDB) Migrate() error {
	return b.db.AutoMigrate(&model.Bookmark{})
}

// GenerateData populates the test database with common user test data.
//
// Returns:
//   - error: An error if data generation fails, otherwise nil
func (b *BookmarkCommonTestDB) GenerateData() error {
	db := b.db.Session(&gorm.Session{})

	bookmarks := []*model.Bookmark{
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "de305d54-75b4-431b-adb2-eb6b9e546000",
			URL:                "https://example.com/alice",
			Description:        "Bookmark for Alice",
			CodeShorten:        1,
			CodeShortenEncoded: "p_1",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000002",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "123e4567-e89b-12d3-a456-eb6b9e546001",
			URL:                "https://example.com/bob",
			Description:        "Bookmark for Bob",
			CodeShorten:        2,
			CodeShortenEncoded: "p_2",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000003",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "987e6543-e21b-12d3-a456-eb6b9e546002",
			URL:                "https://example.com/charlie",
			Description:        "Bookmark for Charlie",
			CodeShorten:        3,
			CodeShortenEncoded: "p_3",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000004",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://example.com/testuser001",
			Description:        "Bookmark for Test User 1 - record 1",
			CodeShorten:        4,
			CodeShortenEncoded: "p_4",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://example.com/testuser001",
			Description:        "Bookmark for Test User 1 - record 2",
			CodeShorten:        5,
			CodeShortenEncoded: "p_5",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
				CreatedAt: TestTime.Add(2 * time.Hour),
				UpdatedAt: TestTime.Add(2 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://example.com/testuser003",
			Description:        "Bookmark for Test User 1 - record 3",
			CodeShorten:        6,
			CodeShortenEncoded: "p_6",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000007",
				CreatedAt: TestTime.Add(3 * time.Hour),
				UpdatedAt: TestTime.Add(3 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://golang.dev/blog/clean-architecture",
			Description:        "Go backend best practices article",
			CodeShorten:        7,
			CodeShortenEncoded: "p_7",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000008",
				CreatedAt: TestTime.Add(4 * time.Hour),
				UpdatedAt: TestTime.Add(4 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://db-tutorials.dev/postgresql-indexing",
			Description:        "Learn PostgreSQL indexing basics",
			CodeShorten:        8,
			CodeShortenEncoded: "p_8",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000009",
				CreatedAt: TestTime.Add(5 * time.Hour),
				UpdatedAt: TestTime.Add(5 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://redis.io/docs/manual/data-types/",
			Description:        "Redis data types documentation",
			CodeShorten:        9,
			CodeShortenEncoded: "p_9",
		},
	}

	return db.CreateInBatches(bookmarks, 10).Error
}
