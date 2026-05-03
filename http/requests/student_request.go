package requests

type StudentPayload struct {
	Name        string `json:"name"`
	ClassroomID uint   `json:"classroom_id"`
}

type StudentUpdatePayload struct {
	Name        *string `json:"name,omitempty"`
	ClassroomID *uint   `json:"classroom_id,omitempty"`
}
