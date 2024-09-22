package request

type NewForum struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	IsActive    *bool   `json:"isActive"`
	PostCount   *uint32 `json:"postCount"`
	CategoryID  *string `json:"categoryId" validate:"required"`
}
