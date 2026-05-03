package resources

import (
	"time"

	"safa-went/database/models"
)

type AuthUserResource struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResource struct {
	Token string           `json:"token"`
	User  AuthUserResource `json:"user"`
}

func NewAuthUserResource(u models.User) AuthUserResource {
	return AuthUserResource{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}