package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type StudentResource struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	ClassroomID uint              `json:"classroom_id"`
	Classroom   ClassroomResource `json:"classroom"`
	Level       string            `json:"level"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type StudentCollection struct {
	Data []StudentResource `json:"data"`
	Meta PaginationMeta    `json:"meta"`
}

type StudentQuery struct {
	db *gorm.DB
}

func NewStudentQuery(db *gorm.DB) *StudentQuery {
	return &StudentQuery{db: db}
}

func NewStudentResource(m models.Student) StudentResource {
	return StudentResource{
		ID:          m.ID,
		Name:        m.Name,
		ClassroomID: m.ClassroomID,
		Classroom:   NewClassroomResource(m.Classroom),
		Level:       m.GetLevel(),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (q *StudentQuery) Paginate(page, perPage int64) (StudentCollection, error) {
	db := q.db.Preload("Classroom")
	data, meta, err := PaginateQuery(db, page, perPage, NewStudentResource)
	if err != nil {
		return StudentCollection{}, err
	}
	return StudentCollection{Data: data, Meta: meta}, nil
}

func (q *StudentQuery) Find(id uint) (StudentResource, error) {
	var item models.Student
	if err := q.db.Preload("Classroom").First(&item, id).Error; err != nil {
		return StudentResource{}, err
	}
	return NewStudentResource(item), nil
}

func (q *StudentQuery) Filter(search string, classroomID uint, page, perPage int64) (StudentCollection, error) {
	db := q.db.Preload("Classroom")
	if search != "" {
		db = db.Where("students.name ILIKE ?", "%"+search+"%")
	}
	if classroomID > 0 {
		db = db.Where("students.classroom_id = ?", classroomID)
	}
	data, meta, err := PaginateQuery(db, page, perPage, NewStudentResource)
	if err != nil {
		return StudentCollection{}, err
	}
	return StudentCollection{Data: data, Meta: meta}, nil
}
