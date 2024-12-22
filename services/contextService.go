package services

import (
	requests "echo-api/models/dtos/requests/context"
	responses "echo-api/models/dtos/responses/pagination"
	"echo-api/models/entities"
	"echo-api/util"
	"errors"
	"fmt"
)

type ContextService struct {
	repo   util.Repository[entities.Context]
	logger *util.Logger
}

func NewContextService(repo util.Repository[entities.Context], logger *util.Logger) *ContextService {
	return &ContextService{repo: repo, logger: logger}
}

func (s *ContextService) CheckIfBelongsToUser(id string, userID string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_CheckIfBelongsToUser with id: %s for user: %s", id, userID))
	context, err := s.repo.First(id, false)
	if err != nil {
		s.logger.Error().Msg("NoteService_CheckIfBelongsToUser had an error when getting from repo")
		return false, err
	}

	return userID == context.UserID, nil
}

func (s *ContextService) GetOne(id string) (entities.Context, error) {
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_GetOne with id: %s", id))

	context, err := s.repo.First(id, true)
	if err != nil {
		s.logger.Error().Msg("ContextService_GetOne had an error when saving to repo")
		return entities.Context{}, err
	}

	return context, nil
}

func (s *ContextService) FilterAll(request requests.FilterContextsRequest) (responses.PaginationResponse[entities.Context], error) {
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	contexts, err := q.Find(true)
	if err != nil {
		s.logger.Error().Msg("ContextService_FilterAll had an error when requesting from repo")
		return responses.PaginationResponse[entities.Context]{}, err
	}
	count, err := q.Count()
	if err != nil {
		s.logger.Error().Msg("ContextService_FilterAll had an error when requesting from repo")
		return responses.PaginationResponse[entities.Context]{}, err
	}
	return responses.PaginationResponse[entities.Context]{Content: contexts, Page: request.Page, Size: len(contexts), TotalCount: int(count)}, nil
}

func (s *ContextService) buildFilterQuery(request requests.FilterContextsRequest) util.Repository[entities.Context] {
	q := s.repo.Query()
	s.logger.Debug().Msg("*ContextService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("*ContextService filtering IDs")
		q = q.Where("ID IN ?", *request.IDs)
	}

	if request.UserIDs != nil && len(*request.UserIDs) > 0 {
		s.logger.Debug().Msg("*ContextService filtering UserIDs")
		q = q.Where("userID IN ?", *request.UserIDs)
	}

	if request.LanguageIDs != nil && len(*request.LanguageIDs) > 0 {
		s.logger.Debug().Msg("*ContextService filtering LanguageIDs")
		q = q.Where("languageID IN ?", *request.LanguageIDs)
	}

	return q.Order("created_at")
}

func (s *ContextService) CreateOne(request requests.CreateContextRequest) (entities.Context, error) {
	if request.LanguageID == "" || request.UserID == "" {
		return entities.Context{}, errors.New("argumentErrorIDMissing")
	}
	context := entities.Context{
		UserID:     request.UserID,
		LanguageID: request.LanguageID,
	}
	s.logger.Debug().Msg("ContextService_CreateOne has started")
	context, err := s.repo.Create(&context)
	if err != nil {
		s.logger.Error().Msg("ContextService_CreateOne had an error when saving to repo")
		return entities.Context{}, err
	}

	return context, nil
}

func (s *ContextService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("ContextService_DeleteOne has started with given id: %s", id))
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error().Msg("ContextService_DeleteOne had an error when deleting from repo")
		return false, err
	}

	return true, nil
}
