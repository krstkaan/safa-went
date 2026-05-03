package models

import (
	"strconv"
	"unicode"

	"gorm.io/gorm"
)

type Classroom struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

// GetLevel derives the school level from the classroom name.
// Names starting with digits 1-4 → "ilkokul", 5-8 → "ortaokul", otherwise → "ortak".
func (c *Classroom) GetLevel() string {
	for _, ch := range c.Name {
		if unicode.IsDigit(ch) {
			n, err := strconv.Atoi(string(ch))
			if err == nil {
				if n >= 1 && n <= 4 {
					return "ilkokul"
				} else if n >= 5 && n <= 8 {
					return "ortaokul"
				}
			}
			break
		}
	}
	return "ortak"
}
