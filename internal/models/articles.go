package models

import "time"

type Article struct {
	ID          string     `gorm:"primaryKey;size:255"`
	Title       string     `gorm:"size:255;not null"`
	Description *string    `gorm:"type:text"`
	Content     *string    `gorm:"type:text"`
	Image       *string    `gorm:"size:255"`
	Author      *string    `gorm:"size:255"`
	Category    *string    `gorm:"size:255;index:idx_posts_category"`
	PublishedAt *time.Time `gorm:"index:idx_posts_published_at"`
	Featured    bool       `gorm:"default:false;index:idx_posts_featured"`
	ReadTime    int        `gorm:"default:0"`
	Excerpt     *string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Article) TableName() string {
	return "articles"
}
