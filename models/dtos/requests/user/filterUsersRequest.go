package user

import (
	base "reson8-learning-api/models/dtos/requests/base"
)

type FilterUsersRequest struct {
	IDs        *[]string `json:"ids" form:"ids"`
	NameQuery  *string   `json:"nameQuery" form:"nameQuery"`
	EmailQuery *string   `json:"emailQuery" form:"emailQuery"`
	CourseIDs  *[]string `json:"courseIds" form:"courseIds"`
	base.PaginationRequestBase
}
