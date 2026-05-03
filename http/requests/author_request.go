package requests

type AuthorPayload struct {
	Name string `json:"name"`
}

type AuthorUpdatePayload struct {
	Name *string `json:"name,omitempty"`
}
