package requests

type BookPayload struct {
	Name        string  `json:"name"`
	Barcode     *string `json:"barcode,omitempty"`
	AuthorID    uint    `json:"author_id"`
	PublisherID uint    `json:"publisher_id"`
	Language    string  `json:"language"`
	PageCount   int     `json:"page_count"`
	IsDonation  bool    `json:"is_donation"`
	ShelfCode   string  `json:"shelf_code"`
	FixtureNo   int     `json:"fixture_no"`
	Level       string  `json:"level"` // ilkokul, ortaokul, ortak
}

type BookUpdatePayload struct {
	Name        *string `json:"name,omitempty"`
	Barcode     *string `json:"barcode,omitempty"`
	AuthorID    *uint   `json:"author_id,omitempty"`
	PublisherID *uint   `json:"publisher_id,omitempty"`
	Language    *string `json:"language,omitempty"`
	PageCount   *int    `json:"page_count,omitempty"`
	IsDonation  *bool   `json:"is_donation,omitempty"`
	ShelfCode   *string `json:"shelf_code,omitempty"`
	FixtureNo   *int    `json:"fixture_no,omitempty"`
	Level       *string `json:"level,omitempty"`
}
