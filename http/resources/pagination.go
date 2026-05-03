package resources

import "gorm.io/gorm"

// PaginationMeta holds pagination metadata in the response envelope.
type PaginationMeta struct {
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`
	PerPage     int64 `json:"per_page"`
	Total       int64 `json:"total"`
	From        int64 `json:"from"`
	To          int64 `json:"to"`
}

// BuildMeta constructs a PaginationMeta from raw pagination values.
func BuildMeta(total, page, perPage, itemCount int64) PaginationMeta {
	lastPage := total / perPage
	if total%perPage != 0 {
		lastPage++
	}

	from := (page-1)*perPage + 1
	to := from + itemCount - 1

	if itemCount == 0 {
		from = 0
		to = 0
	}

	return PaginationMeta{
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     perPage,
		Total:       total,
		From:        from,
		To:          to,
	}
}

// PaginateQuery is a generic helper that handles clamp, COUNT, offset/limit,
// and transform. Pass a pre-scoped db (e.g. with Where clauses already applied).
func PaginateQuery[M any, R any](
	db *gorm.DB,
	page, perPage int64,
	transform func(M) R,
) ([]R, PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}
	if perPage > 100 {
		perPage = 100
	}

	var model M
	var total int64
	if err := db.Model(&model).Count(&total).Error; err != nil {
		return nil, PaginationMeta{}, err
	}

	var items []M
	offset := (page - 1) * perPage
	if err := db.Offset(int(offset)).Limit(int(perPage)).Find(&items).Error; err != nil {
		return nil, PaginationMeta{}, err
	}

	data := make([]R, len(items))
	for i, item := range items {
		data[i] = transform(item)
	}

	return data, BuildMeta(total, page, perPage, int64(len(items))), nil
}
