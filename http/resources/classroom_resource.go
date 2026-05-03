package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type ClassroomResource struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClassroomCollection struct {
	Data []ClassroomResource `json:"data"`
	Meta PaginationMeta      `json:"meta"`
}

type ClassroomQuery struct {
	db *gorm.DB
}

func NewClassroomQuery(db *gorm.DB) *ClassroomQuery {
	return &ClassroomQuery{db: db}
}

func NewClassroomResource(m models.Classroom) ClassroomResource {
	return ClassroomResource{
		ID:        m.ID,
		Name:      m.Name,
		Level:     m.GetLevel(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (q *ClassroomQuery) Paginate(page, perPage int64) (ClassroomCollection, error) {
	data, meta, err := PaginateQuery(q.db, page, perPage, NewClassroomResource)
	if err != nil {
		return ClassroomCollection{}, err
	}
	return ClassroomCollection{Data: data, Meta: meta}, nil
}

func (q *ClassroomQuery) Find(id uint) (ClassroomResource, error) {
	var item models.Classroom
	if err := q.db.First(&item, id).Error; err != nil {
		return ClassroomResource{}, err
	}
	return NewClassroomResource(item), nil
}

func (q *ClassroomQuery) Search(name string, page, perPage int64) (ClassroomCollection, error) {
	filtered := q.db.Where("name ILIKE ?", "%"+name+"%")
	data, meta, err := PaginateQuery(filtered, page, perPage, NewClassroomResource)
	if err != nil {
		return ClassroomCollection{}, err
	}
	return ClassroomCollection{Data: data, Meta: meta}, nil
}
