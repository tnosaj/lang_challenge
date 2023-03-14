package domain

type Order struct {
	ID     string `json:"ID" validate:"required"`
	Status string `json:"Status" validate:"required"`
}
