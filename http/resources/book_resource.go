package resources

import (
	"time"

	"gorm.io/gorm"

	"safa-went/database/models"
)

type BookResource struct {
	ID                uint              `json:"id"`
	Name              string            `json:"name"`
	Barcode           *string           `json:"barcode"`
	AuthorID          uint              `json:"author_id"`
	Author            AuthorResource    `json:"author"`
	PublisherID       uint              `json:"publisher_id"`
	Publisher         PublisherResource `json:"publisher"`
	Language          string            `json:"language"`
	PageCount         int               `json:"page_count"`
	IsDonation        bool              `json:"is_donation"`
	ShelfCode         string            `json:"shelf_code"`
	FixtureNo         int               `json:"fixture_no"`
	Level             string            `json:"level"`
	IsCurrentlyLoaned bool              `json:"is_currently_loaned"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type BookCollection struct {
	Data []BookResource `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type BookQuery struct {
	db *gorm.DB
}

func NewBookQuery(db *gorm.DB) *BookQuery {
	return &BookQuery{db: db}
}

func newBookResourceWithLoan(m models.Book, db *gorm.DB) BookResource {
	var loanCount int64
	db.Model(&models.Loan{}).
		Where("book_id = ? AND status = ? AND deleted_at IS NULL", m.ID, models.LoanStatusActive).
		Count(&loanCount)

	return BookResource{
		ID:                m.ID,
		Name:              m.Name,
		Barcode:           m.Barcode,
		AuthorID:          m.AuthorID,
		Author:            NewAuthorResource(m.Author),
		PublisherID:       m.PublisherID,
		Publisher:         NewPublisherResource(m.Publisher),
		Language:          m.Language,
		PageCount:         m.PageCount,
		IsDonation:        m.IsDonation,
		ShelfCode:         m.ShelfCode,
		FixtureNo:         m.FixtureNo,
		Level:             m.Level,
		IsCurrentlyLoaned: loanCount > 0,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func (q *BookQuery) Find(id uint) (BookResource, error) {
	var item models.Book
	if err := q.db.Preload("Author").Preload("Publisher").First(&item, id).Error; err != nil {
		return BookResource{}, err
	}
	return newBookResourceWithLoan(item, q.db), nil
}

// Filter applies optional filters and returns a paginated collection.
func (q *BookQuery) Filter(
	name string,
	fixtureNo int,
	search string,
	authorID uint,
	publisherID uint,
	level string,
	isDonation *bool,
	sortBy string,
	page, perPage int64,
) (BookCollection, error) {
	db := q.db.Preload("Author").Preload("Publisher")
	if name != "" {
		db = db.Where("books.name ILIKE ?", "%"+name+"%")
	}
	if fixtureNo > 0 {
		db = db.Where("books.fixture_no = ?", fixtureNo)
	}
	if search != "" {
		db = db.Where("books.name ILIKE ? OR books.barcode ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if authorID > 0 {
		db = db.Where("books.author_id = ?", authorID)
	}
	if publisherID > 0 {
		db = db.Where("books.publisher_id = ?", publisherID)
	}
	if level != "" {
		db = db.Where("books.level = ?", level)
	}
	if isDonation != nil {
		db = db.Where("books.is_donation = ?", *isDonation)
	}
	if sortBy != "" {
		db = db.Order(sortBy)
	} else {
		db = db.Order("books.fixture_no ASC")
	}

	// Use a slice-based paginate since we need per-item DB calls for IsCurrentlyLoaned
	var items []models.Book
	var total int64
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Model(&models.Book{}).Count(&total).Error; err != nil {
		return BookCollection{}, err
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage
	if err := db.Offset(int(offset)).Limit(int(perPage)).Find(&items).Error; err != nil {
		return BookCollection{}, err
	}

	data := make([]BookResource, 0, len(items))
	for _, item := range items {
		data = append(data, newBookResourceWithLoan(item, q.db))
	}

	meta := BuildMeta(total, page, perPage, int64(len(items)))
	return BookCollection{Data: data, Meta: meta}, nil
}
