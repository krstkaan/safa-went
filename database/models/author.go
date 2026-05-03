package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name string `json:"name" gorm:"not null;uniqueIndex"`
}
