package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SelectedVariant struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type SelectedVariants []SelectedVariant

func (s SelectedVariants) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (s *SelectedVariants) Scan(src any) error {
	if src == nil {
		*s = SelectedVariants{}
		return nil
	}

	var raw []byte
	switch v := src.(type) {
	case string:
		raw = []byte(v)
	case []byte:
		raw = v
	default:
		return fmt.Errorf("unsupported scan type for SelectedVariants: %T", src)
	}

	if len(raw) == 0 {
		*s = SelectedVariants{}
		return nil
	}

	return json.Unmarshal(raw, s)
}

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uint      `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time

	Items []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}

type CartItem struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CartID           uuid.UUID        `gorm:"type:uuid;not null;index"`
	ProductID        uuid.UUID        `gorm:"type:uuid;not null;index"`
	Quantity         int              `gorm:"not null"`
	SelectedVariants SelectedVariants `gorm:"type:jsonb;not null;default:'[]'"`
	CreatedAt        time.Time        `gorm:"<-:create"`
	UpdatedAt        time.Time

	Product Product `gorm:"foreignKey:ProductID"`
}

func (Cart) TableName() string     { return "carts" }
func (CartItem) TableName() string { return "cart_items" }
