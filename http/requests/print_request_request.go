package requests

import "time"

type PrintRequestPayload struct {
	RequestedAt time.Time `json:"requested_at"`
	ColorCopies int       `json:"color_copies"`
	BWCopies    int       `json:"bw_copies"`
	Description *string   `json:"description,omitempty"`
	RequesterID uint      `json:"requester_id"`
	ApproverID  uint      `json:"approver_id"`
}

type PrintRequestUpdatePayload struct {
	RequestedAt *time.Time `json:"requested_at,omitempty"`
	ColorCopies *int       `json:"color_copies,omitempty"`
	BWCopies    *int       `json:"bw_copies,omitempty"`
	Description *string    `json:"description,omitempty"`
	RequesterID *uint      `json:"requester_id,omitempty"`
	ApproverID  *uint      `json:"approver_id,omitempty"`
}
