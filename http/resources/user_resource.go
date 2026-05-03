package resources

import (
	"safa-went/database/models"
	"time"

	"gorm.io/gorm"
)

// UserResource is the single-item representation returned by the API.
type UserResource struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCollection is the paginated list representation returned by the API.
type UserCollection struct {
	Data []UserResource `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

// UserQuery handles all DB read operations for User and transforms results.
type UserQuery struct {
	db *gorm.DB
}

// NewUserQuery creates a new UserQuery bound to the given DB instance.
func NewUserQuery(db *gorm.DB) *UserQuery {
	return &UserQuery{db: db}
}

// Paginate runs a COUNT + paginated SELECT and returns a UserCollection.
func (q *UserQuery) Paginate(page, perPage int64) (UserCollection, error) {
	data, meta, err := PaginateQuery(q.db, page, perPage, NewUserResource)
	if err != nil {
		return UserCollection{}, err
	}
	return UserCollection{Data: data, Meta: meta}, nil
}

// Find fetches a single User by primary key and returns a UserResource.
func (q *UserQuery) Find(id uint) (UserResource, error) {
	var m models.User
	if err := q.db.First(&m, id).Error; err != nil {
		return UserResource{}, err
	}
	return NewUserResource(m), nil
}

// NewUserResource transforms a models.User into a UserResource.
func NewUserResource(m models.User) UserResource {
	return UserResource{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
