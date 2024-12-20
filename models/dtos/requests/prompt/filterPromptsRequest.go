package prompt

import (
	"reson8-learning-api/models/dtos/requests/base"
)

type FilterPromptsRequest struct {
	base.PaginationRequestBase
	IDs        *[]string `json:"ids" form:"ids"`
	ContextIDs *[]string `json:"contexts" form:"contexts"`
}
