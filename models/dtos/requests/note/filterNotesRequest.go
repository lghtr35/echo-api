package note

import (
	base "echo-api/models/dtos/requests/base"
)

type FilterNotesRequest struct {
	IDs         *[]string `json:"ids" form:"ids"`
	Header      *string   `json:"name" form:"name"`
	UserIDs     *[]string `json:"users" form:"users"`
	DocumentIDs *[]string `json:"documents" form:"documents"`
	LanguageIDs *[]string `json:"languages" form:"languages"`
	ContextIDs  *[]string `json:"contexts" form:"contexts"`
	base.PaginationRequestBase
}
