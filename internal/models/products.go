package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string      `gorm:"size:255;not null"`
	Description   string      `gorm:"type:text;not null;default:''"`
	CategoryID    uint        `gorm:"not null"`
	Price         *float64    `gorm:"type:decimal(15,2);not null"`
	OriginalPrice *float64    `gorm:"type:decimal(15,2)"`
	Image         string      `gorm:"type:text;not null"`
	Rating        float64     `gorm:"type:decimal(2,1);not null;default:0"`
	Reviews       int         `gorm:"not null;default:0"`
	Stock         int         `gorm:"not null;default:0"`
	InStock       bool        `gorm:"not null;default:true"`
	Featured      bool        `gorm:"not null;default:false"`
	Tags          StringSlice `gorm:"type:text[];not null;default:'{}'"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Category Category         `gorm:"foreignKey:CategoryID"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type ProductVariant struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Name            string    `gorm:"size:100;not null"`
	Type            string    `gorm:"size:50;not null"`
	Value           string    `gorm:"size:100;not null"`
	Price           *float64  `gorm:"type:decimal(15,2)"`
	PriceAdjustment *float64  `gorm:"type:decimal(15,2)"`
	Stock           int       `gorm:"not null;default:0"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

func (ProductVariant) TableName() string {
	return "product_variants"
}
