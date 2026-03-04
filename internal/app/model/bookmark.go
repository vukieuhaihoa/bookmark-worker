package model

// Bookmark represents a bookmark entity in the system.
type Bookmark struct {
	Base

	Description string `gorm:"column:description" json:"description"`
	URL         string `gorm:"not null;column:url" json:"url"`

	UserID string `gorm:"not null;column:user_id" json:"-"`

	CodeShorten        int64  `gorm:"column:code_shorten;autoIncrement" json:"-"`
	CodeShortenEncoded string `gorm:"column:code_shorten_encoded" json:"code"`
}
