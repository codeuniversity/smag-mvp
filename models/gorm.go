package models

import "time"

// GormModelWithoutID is an alternative definition for gorm.Model without an ID
type GormModelWithoutID struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
