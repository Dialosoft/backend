package request

type CreateForum struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type" validate:"required"`
}

type UpdateForum struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}
