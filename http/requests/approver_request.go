package requests

type ApproverPayload struct {
	Name string `json:"name"`
}

type ApproverUpdatePayload struct {
	Name *string `json:"name,omitempty"`
}
