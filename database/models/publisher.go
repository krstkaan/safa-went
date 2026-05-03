package models

import "gorm.io/gorm"

type Publisher struct {
	gorm.Model
	Name string `json:"name" gorm:"not null;uniqueIndex"`
}
