package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"size:50;not null"`
	Slug     string    `gorm:"size:50;uniqueIndex;not null"`
	ParentID *uint     `gorm:"index" json:"parent_id"`
	Parent   *Category `gorm:"foreignKey:ParentID" json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Products []Product
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
