package resources

import "gorm.io/gorm"

// auth_controllerResource is a skeleton resource generated without --all.
type auth_controllerResource struct {}

// auth_controllerCollection is a skeleton collection generated without --all.
type auth_controllerCollection struct {
	Data []auth_controllerResource `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

// auth_controllerQuery is a skeleton query builder generated without --all.
type auth_controllerQuery struct {
	db *gorm.DB
}

func Newauth_controllerQuery(db *gorm.DB) *auth_controllerQuery {
	return &auth_controllerQuery{db: db}
}
