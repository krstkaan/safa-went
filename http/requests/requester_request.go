package requests

type RequesterPayload struct {
	Name string `json:"name"`
}

type RequesterUpdatePayload struct {
	Name *string `json:"name,omitempty"`
}
