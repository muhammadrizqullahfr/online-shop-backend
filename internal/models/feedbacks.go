package models

import (
	"time"

	"gorm.io/gorm"
)

type Feedback struct {
	ID        int            `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"size:255;not null"`
	Email     string         `gorm:"size:255;not null"`
	Subject   string         `gorm:"size:255;not null"`
	Message   string         `gorm:"type:text;not null;default:''"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}
