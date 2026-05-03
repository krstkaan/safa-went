package models

import "gorm.io/gorm"

type Approver struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}
