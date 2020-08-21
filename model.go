package lazy

import (
	"database/sql"
	"time"
)

// DeletedAt compatible with gorm
type DeletedAt sql.NullTime

// Model compatible with gorm, json, struct
type Model struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
}

// ModelDelete compatible with gorm, json, struct and DeletedAt.
type ModelDelete struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
	DeletedAt DeletedAt `gorm:"index" json:"deleted_at" lazy:"deleted_at" mapstructure:"deleted_at"`
}
