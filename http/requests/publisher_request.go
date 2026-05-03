package requests

type PublisherPayload struct {
	Name string `json:"name"`
}

type PublisherUpdatePayload struct {
	Name *string `json:"name,omitempty"`
}
