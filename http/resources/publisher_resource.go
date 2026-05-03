package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type PublisherResource struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublisherCollection struct {
	Data []PublisherResource `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}

type PublisherQuery struct {
	db *gorm.DB
}

func NewPublisherQuery(db *gorm.DB) *PublisherQuery {
	return &PublisherQuery{db: db}
}

func NewPublisherResource(m models.Publisher) PublisherResource {
	return PublisherResource{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (q *PublisherQuery) Paginate(page, perPage int64) (PublisherCollection, error) {
	data, meta, err := PaginateQuery(q.db, page, perPage, NewPublisherResource)
	if err != nil {
		return PublisherCollection{}, err
	}
	return PublisherCollection{Data: data, Meta: meta}, nil
}

func (q *PublisherQuery) Find(id uint) (PublisherResource, error) {
	var item models.Publisher
	if err := q.db.First(&item, id).Error; err != nil {
		return PublisherResource{}, err
	}
	return NewPublisherResource(item), nil
}

func (q *PublisherQuery) Search(name string, page, perPage int64) (PublisherCollection, error) {
	filtered := q.db.Where("name ILIKE ?", "%"+name+"%")
	data, meta, err := PaginateQuery(filtered, page, perPage, NewPublisherResource)
	if err != nil {
		return PublisherCollection{}, err
	}
	return PublisherCollection{Data: data, Meta: meta}, nil
}
