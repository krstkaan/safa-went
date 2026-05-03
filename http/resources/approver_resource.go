package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type ApproverResource struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ApproverCollection struct {
	Data []ApproverResource `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

type ApproverQuery struct {
	db *gorm.DB
}

func NewApproverQuery(db *gorm.DB) *ApproverQuery {
	return &ApproverQuery{db: db}
}

func NewApproverResource(m models.Approver) ApproverResource {
	return ApproverResource{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (q *ApproverQuery) Paginate(page, perPage int64) (ApproverCollection, error) {
	data, meta, err := PaginateQuery(q.db, page, perPage, NewApproverResource)
	if err != nil {
		return ApproverCollection{}, err
	}
	return ApproverCollection{Data: data, Meta: meta}, nil
}

func (q *ApproverQuery) Find(id uint) (ApproverResource, error) {
	var item models.Approver
	if err := q.db.First(&item, id).Error; err != nil {
		return ApproverResource{}, err
	}
	return NewApproverResource(item), nil
}

func (q *ApproverQuery) Search(name string, page, perPage int64) (ApproverCollection, error) {
	filtered := q.db.Where("name ILIKE ?", "%"+name+"%")
	data, meta, err := PaginateQuery(filtered, page, perPage, NewApproverResource)
	if err != nil {
		return ApproverCollection{}, err
	}
	return ApproverCollection{Data: data, Meta: meta}, nil
}
