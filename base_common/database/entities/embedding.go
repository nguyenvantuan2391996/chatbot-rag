package entities

import (
	"time"

	"gorm.io/gorm"
)

type Embedding struct {
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Vector    string         `json:"vector"`
	ID        int64          `json:"id"`
}
