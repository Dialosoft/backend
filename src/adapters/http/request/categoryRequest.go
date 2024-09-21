package request

type CreateCategory struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateCategory struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
