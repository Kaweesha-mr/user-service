package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username string    `gorm:"not null" validate:"required,min=3,max=100"`
	Email    string    `gorm:"unique;not null" validate:"required,email"`
}
