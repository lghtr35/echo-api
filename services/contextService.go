package services

import (
	"errors"
	"fmt"
	requests "reson8-learning-api/models/dtos/requests/context"
	responses "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContextService struct {
	db     *gorm.DB
	logger *util.Logger
}

func NewContextService(db *gorm.DB, logger *util.Logger) *ContextService {
	return &ContextService{db: db, logger: logger}
}

func (s *ContextService) CheckIfBelongsToUser(id string, userID string) (bool, error) {
	var dbID string
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_CheckIfBelongsToUser with id: %s for user: %s", id, userID))
	res := s.db.Model(entities.Context{}).Pluck("userID", &dbID)
	if res.Error != nil {
		s.logger.Error().Msg("ContextService_CheckIfBelongsToUser had an error when getting from db")
		return false, res.Error
	}

	return userID == dbID, nil
}

func (s ContextService) GetOne(id string) (entities.Context, error) {
	var context entities.Context
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_GetOne with id: %s", id))

	res := s.db.First(&context, id)
	if res.Error != nil {
		s.logger.Error().Msg("ContextService_GetOne had an error when saving to db")
		return entities.Context{}, res.Error
	}

	return context, nil
}

func (s ContextService) FilterAll(request requests.FilterContextsRequest) (responses.PaginationResponse[entities.Context], error) {
	var notes []entities.Context
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	res := q.Find(&notes)
	if res.Error != nil {
		s.logger.Error().Msg("ContextService_FilterAll had an error when requesting from db")
		return responses.PaginationResponse[entities.Context]{}, res.Error
	}
	var count int64
	res = q.Count(&count)
	if res.Error != nil {
		s.logger.Error().Msg("ContextService_FilterAll had an error when requesting from db")
		return responses.PaginationResponse[entities.Context]{}, res.Error
	}
	return responses.PaginationResponse[entities.Context]{Content: notes, Page: request.Page, Size: len(notes), TotalCount: int(count)}, nil
}

func (s ContextService) buildFilterQuery(request requests.FilterContextsRequest) *gorm.DB {
	q := s.db.Model(&entities.Context{})
	s.logger.Debug().Msg("ContextService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("ContextService filtering IDs")
		q = q.Where(*request.IDs)
	}

	if request.UserIDs != nil && len(*request.UserIDs) > 0 {
		s.logger.Debug().Msg("ContextService filtering UserIDs")
		q = q.Where("userID IN ?", *request.UserIDs)
	}

	if request.LanguageIDs != nil && len(*request.LanguageIDs) > 0 {
		s.logger.Debug().Msg("ContextService filtering LanguageIDs")
		q = q.Where("languageID IN ?", *request.LanguageIDs)
	}

	return q.Order("created_at")
}

func (s ContextService) CreateOne(request requests.CreateContextRequest) (entities.Context, error) {
	if request.LanguageID == "" || request.UserID == "" {
		return entities.Context{}, errors.New("argumentErrorIDMissing")
	}
	context := entities.Context{
		UserID:     request.UserID,
		LanguageID: request.LanguageID,
	}
	s.logger.Debug().Msg("ContextService_CreateOne has started")
	res := s.db.Create(&context)
	if res.Error != nil {
		s.logger.Error().Msg("ContextService_CreateOne had an error when saving to db")
		return entities.Context{}, res.Error
	}

	return context, nil
}

func (s ContextService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_DeleteOne has started with given id: %s", id))
	res := s.db.Select(clause.Associations).Delete(&entities.Context{}, id)
	if res.Error != nil {
		s.logger.Error().Msg("ContextService_DeleteOne had an error when deleting from db")
		return false, res.Error
	}

	return true, nil
}
