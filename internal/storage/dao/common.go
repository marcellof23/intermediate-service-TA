package dao

import "time"

type Base struct {
	CreatedAt time.Time `gorm:"index:idx_created_at;not null"`
	UpdatedAt time.Time `gorm:"index:idx_updated_at;not null"`
}
