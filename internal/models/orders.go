package models

import "time"

type Order struct {
	ID              string  `gorm:"type:varchar(255);primaryKey"`
	CustomerName    string  `gorm:"size:255;not null"`
	CustomerEmail   string  `gorm:"size:255;not null"`
	CustomerPhone   string  `gorm:"size:255;not null"`
	Total           float64 `gorm:"type:float8;not null"`
	Status          *string `gorm:"size:255"`
	PaymentURL      *string `gorm:"size:255;column:payment_url"`
	UserID          *string `gorm:"size:255"`
	CustomerAddress string  `gorm:"type:text;not null"`
	City            string  `gorm:"size:255;not null"`
	State           string  `gorm:"size:255;not null"`
	ZipCode         string  `gorm:"size:255;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Items []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	ID               string           `gorm:"type:varchar(255);primaryKey"`
	OrderID          string           `gorm:"type:varchar(255);not null;index"`
	ProductID        string           `gorm:"type:varchar(255);not null;index"`
	Name             string           `gorm:"size:255;not null"`
	Price            float64          `gorm:"type:float8;not null"`
	Quantity         int              `gorm:"not null"`
	Image            *string          `gorm:"size:255"`
	SelectedVariants SelectedVariants `gorm:"type:jsonb"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (Order) TableName() string     { return "orders" }
func (OrderItem) TableName() string { return "order_items" }
