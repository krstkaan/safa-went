package models

import "gorm.io/gorm"

type Requester struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}
