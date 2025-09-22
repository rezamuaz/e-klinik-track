package dto

type SuccessResponse struct {
	Data       any `json:"data,omitempty"`
	Pagination any `json:"pagination,omitempty"`
}
