package context

import "echo-api/models/dtos/requests/base"

type FilterContextsRequest struct {
	base.PaginationRequestBase
	IDs         *[]string `json:"ids" form:"ids"`
	UserIDs     *[]string `json:"userIds" form:"userIds"`
	LanguageIDs *[]string `json:"languageIds" form:"languageIds"`
}
