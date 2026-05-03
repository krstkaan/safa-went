package models

import (
	"time"

	"gorm.io/gorm"
)

type PrintRequest struct {
	gorm.Model
	RequestedAt time.Time `json:"requested_at" gorm:"not null"`
	ColorCopies int       `json:"color_copies" gorm:"not null;default:0"`
	BWCopies    int       `json:"bw_copies" gorm:"not null;default:0"`
	Description *string   `json:"description"`
	RequesterID uint      `json:"requester_id" gorm:"not null"`
	Requester   Requester `json:"requester" gorm:"foreignKey:RequesterID"`
	ApproverID  uint      `json:"approver_id" gorm:"not null"`
	Approver    Approver  `json:"approver" gorm:"foreignKey:ApproverID"`
}
