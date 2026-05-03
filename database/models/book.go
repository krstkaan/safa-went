package models

import "gorm.io/gorm"

// BookLevel enum values.
const (
	BookLevelIlkokul  = "ilkokul"
	BookLevelOrtaokul = "ortaokul"
	BookLevelOrtak    = "ortak"
)

type Book struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Barcode     *string   `json:"barcode" gorm:"uniqueIndex"`
	AuthorID    uint      `json:"author_id" gorm:"not null"`
	Author      Author    `json:"author" gorm:"foreignKey:AuthorID"`
	PublisherID uint      `json:"publisher_id" gorm:"not null"`
	Publisher   Publisher `json:"publisher" gorm:"foreignKey:PublisherID"`
	Language    string    `json:"language" gorm:"not null"`
	PageCount   int       `json:"page_count" gorm:"not null"`
	IsDonation  bool      `json:"is_donation" gorm:"default:false"`
	ShelfCode   string    `json:"shelf_code" gorm:"not null"`
	FixtureNo   int       `json:"fixture_no" gorm:"uniqueIndex;not null"`
	Level       string    `json:"level" gorm:"type:varchar(20);not null"` // ilkokul, ortaokul, ortak
}
