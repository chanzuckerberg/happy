package model

import (
	"time"

	"gorm.io/gorm"
)

type CommonDBFields struct {
	ID        uint           `gorm:"primarykey" swaggerignore:"true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
}
