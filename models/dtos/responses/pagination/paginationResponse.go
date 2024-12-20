package pagination

// PaginationResponse[T]
// @Description Pagination response
type PaginationResponse[T any] struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalCount int `json:"totalCount"`
	Content    []T `json:"content"`
}
