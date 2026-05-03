package resources

import "gorm.io/gorm"

// authResource is a skeleton resource generated without --all.
type authResource struct {}

// authCollection is a skeleton collection generated without --all.
type authCollection struct {
	Data []authResource `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

// authQuery is a skeleton query builder generated without --all.
type authQuery struct {
	db *gorm.DB
}

func NewauthQuery(db *gorm.DB) *authQuery {
	return &authQuery{db: db}
}
