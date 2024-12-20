package base

type PaginationRequestBase struct {
	Page int `json:"page" form:"page" binding:"required"`
	Size int `json:"size" form:"size" binding:"required"`
}

func (prb *PaginationRequestBase) CalculateOffset() int {
	return (prb.Page - 1) * prb.Size
}
