package document

import (
	base_request "reson8-learning-api/models/dtos/requests/base"
)

type FilterDocumentsRequest struct {
	IDs        *[]string `json:"ids" form:"ids"`
	UserIDs    *[]string `json:"userIds" form:"userIds"`
	Name       *string   `json:"name" form:"name"`
	Location   *string   `json:"location" form:"location"`
	Extension  *string   `json:"extension" form:"extension"`
	NoteIDs    *[]string `json:"noteIds" form:"noteIds"`
	ContextIDs *[]string `json:"contextIds" form:"contextIds"`
	base_request.PaginationRequestBase
}
