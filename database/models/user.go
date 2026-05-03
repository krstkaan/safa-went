package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"-" gorm:"not null"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
}
