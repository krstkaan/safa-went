package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type RequesterResource struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequesterCollection struct {
	Data []RequesterResource `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}

type RequesterQuery struct {
	db *gorm.DB
}

func NewRequesterQuery(db *gorm.DB) *RequesterQuery {
	return &RequesterQuery{db: db}
}

func NewRequesterResource(m models.Requester) RequesterResource {
	return RequesterResource{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (q *RequesterQuery) Paginate(page, perPage int64) (RequesterCollection, error) {
	data, meta, err := PaginateQuery(q.db, page, perPage, NewRequesterResource)
	if err != nil {
		return RequesterCollection{}, err
	}
	return RequesterCollection{Data: data, Meta: meta}, nil
}

func (q *RequesterQuery) Find(id uint) (RequesterResource, error) {
	var item models.Requester
	if err := q.db.First(&item, id).Error; err != nil {
		return RequesterResource{}, err
	}
	return NewRequesterResource(item), nil
}

func (q *RequesterQuery) Search(name string, page, perPage int64) (RequesterCollection, error) {
	filtered := q.db.Where("name ILIKE ?", "%"+name+"%")
	data, meta, err := PaginateQuery(filtered, page, perPage, NewRequesterResource)
	if err != nil {
		return RequesterCollection{}, err
	}
	return RequesterCollection{Data: data, Meta: meta}, nil
}
