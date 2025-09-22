package dto

type Pagination struct {
	CurrentPage int64 `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
}
