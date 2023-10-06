package model

import (
	"time"

	"database/sql"
)

type CommonDBFields struct {
	ID        uint         `json:"id" swaggerignore:"true"`
	CreatedAt time.Time    `json:"created_at" example:"2023-02-22T22:55:02.20034Z"`
	UpdatedAt time.Time    `json:"updated_at" example:"2023-02-22T22:55:02.20034Z"`
	DeletedAt sql.NullTime `json:"deleted_at" swaggerignore:"true"`
}
