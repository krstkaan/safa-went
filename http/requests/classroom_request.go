package requests

type ClassroomPayload struct {
	Name string `json:"name"`
}

type ClassroomUpdatePayload struct {
	Name *string `json:"name,omitempty"`
}
