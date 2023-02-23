package model

import (
	"time"

	"gorm.io/gorm"
)

type CommonDBFields struct {
	ID        uint           `gorm:"primarykey" swaggerignore:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2023-02-22T22:55:02.20034Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-02-22T22:55:02.20034Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
}
