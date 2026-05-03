package requests

// UserPayload represents required fields for create operations.
type UserPayload struct {
	Email string `json:"email" example:"john.doe@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

// UserUpdatePayload represents optional fields for partial updates.
type UserUpdatePayload struct {
	Email *string `json:"email,omitempty" example:"john.doe@example.com"`
	Name  *string `json:"name,omitempty" example:"John Doe"`
}
