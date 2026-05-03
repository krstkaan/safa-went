package resources

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
