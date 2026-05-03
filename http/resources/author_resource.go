package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type AuthorResource struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthorCollection struct {
	Data []AuthorResource `json:"data"`
	Meta PaginationMeta   `json:"meta"`
}

type AuthorQuery struct {
	db *gorm.DB
}

func NewAuthorQuery(db *gorm.DB) *AuthorQuery {
	return &AuthorQuery{db: db}
}

func NewAuthorResource(m models.Author) AuthorResource {
	return AuthorResource{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (q *AuthorQuery) Paginate(page, perPage int64) (AuthorCollection, error) {
	data, meta, err := PaginateQuery(q.db, page, perPage, NewAuthorResource)
	if err != nil {
		return AuthorCollection{}, err
	}
	return AuthorCollection{Data: data, Meta: meta}, nil
}

func (q *AuthorQuery) Find(id uint) (AuthorResource, error) {
	var item models.Author
	if err := q.db.First(&item, id).Error; err != nil {
		return AuthorResource{}, err
	}
	return NewAuthorResource(item), nil
}

func (q *AuthorQuery) Search(name string, page, perPage int64) (AuthorCollection, error) {
	filtered := q.db.Where("name ILIKE ?", "%"+name+"%")
	data, meta, err := PaginateQuery(filtered, page, perPage, NewAuthorResource)
	if err != nil {
		return AuthorCollection{}, err
	}
	return AuthorCollection{Data: data, Meta: meta}, nil
}
