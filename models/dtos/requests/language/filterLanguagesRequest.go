package language

import (
	base "reson8-learning-api/models/dtos/requests/base"
)

type FilterLanguagesRequest struct {
	IDs        *[]string `json:"ids" form:"ids"`
	Name       *string   `json:"name" form:"name"`
	Alpha2Code *string   `json:"alpha2Code" form:"alpha2Code"`
	UserIDs    *[]string `json:"userIds" form:"userIds"`
	CourseIDs  *[]string `json:"courseIds" form:"courseIds"`
	base.PaginationRequestBase
}
