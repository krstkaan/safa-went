package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type PrintRequestResource struct {
	ID          uint              `json:"id"`
	RequestedAt time.Time         `json:"requested_at"`
	ColorCopies int               `json:"color_copies"`
	BWCopies    int               `json:"bw_copies"`
	Description *string           `json:"description"`
	RequesterID uint              `json:"requester_id"`
	Requester   RequesterResource `json:"requester"`
	ApproverID  uint              `json:"approver_id"`
	Approver    ApproverResource  `json:"approver"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type PrintRequestCollection struct {
	Data []PrintRequestResource `json:"data"`
	Meta PaginationMeta         `json:"meta"`
}

type PrintRequestQuery struct {
	db *gorm.DB
}

func NewPrintRequestQuery(db *gorm.DB) *PrintRequestQuery {
	return &PrintRequestQuery{db: db}
}

func NewPrintRequestResource(m models.PrintRequest) PrintRequestResource {
	return PrintRequestResource{
		ID:          m.ID,
		RequestedAt: m.RequestedAt,
		ColorCopies: m.ColorCopies,
		BWCopies:    m.BWCopies,
		Description: m.Description,
		RequesterID: m.RequesterID,
		Requester:   NewRequesterResource(m.Requester),
		ApproverID:  m.ApproverID,
		Approver:    NewApproverResource(m.Approver),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// Filter applies optional filters and returns a paginated collection.
func (q *PrintRequestQuery) Filter(
	requesterNames []string,
	approverNames []string,
	colorCopiesMin, colorCopiesMax int,
	bwCopiesMin, bwCopiesMax int,
	requestedAtFrom, requestedAtTo *time.Time,
	sortBy string,
	page, perPage int64,
) (PrintRequestCollection, error) {
	db := q.db.Preload("Requester").Preload("Approver")
	if len(requesterNames) > 0 {
		db = db.Joins("JOIN requesters ON requesters.id = print_requests.requester_id").
			Where("requesters.name IN ?", requesterNames)
	}
	if len(approverNames) > 0 {
		db = db.Joins("JOIN approvers ON approvers.id = print_requests.approver_id").
			Where("approvers.name IN ?", approverNames)
	}
	if colorCopiesMin > 0 {
		db = db.Where("print_requests.color_copies >= ?", colorCopiesMin)
	}
	if colorCopiesMax > 0 {
		db = db.Where("print_requests.color_copies <= ?", colorCopiesMax)
	}
	if bwCopiesMin > 0 {
		db = db.Where("print_requests.bw_copies >= ?", bwCopiesMin)
	}
	if bwCopiesMax > 0 {
		db = db.Where("print_requests.bw_copies <= ?", bwCopiesMax)
	}
	if requestedAtFrom != nil {
		db = db.Where("print_requests.requested_at >= ?", *requestedAtFrom)
	}
	if requestedAtTo != nil {
		db = db.Where("print_requests.requested_at <= ?", *requestedAtTo)
	}
	if sortBy != "" {
		db = db.Order(sortBy)
	} else {
		db = db.Order("print_requests.requested_at DESC")
	}

	data, meta, err := PaginateQuery(db, page, perPage, NewPrintRequestResource)
	if err != nil {
		return PrintRequestCollection{}, err
	}
	return PrintRequestCollection{Data: data, Meta: meta}, nil
}

func (q *PrintRequestQuery) Find(id uint) (PrintRequestResource, error) {
	var item models.PrintRequest
	if err := q.db.Preload("Requester").Preload("Approver").First(&item, id).Error; err != nil {
		return PrintRequestResource{}, err
	}
	return NewPrintRequestResource(item), nil
}
