package user

import (
	"reson8-learning-api/models/dtos/operation"
	"reson8-learning-api/models/entities"
)

type UpdateUserRequest struct {
	ID        string                                   `json:"id"`
	Name      *string                                  `json:"name"`
	Notes     *[]operation.Operable[entities.Note]     `json:"noteOps"`
	Languages *[]operation.Operable[entities.Language] `json:"languageOps"`
}

type Role uint

const (
	Student Role = iota
	Teacher
	StudentTeacher
)
