package services

import (
	requests "echo-api/models/dtos/requests/language"
	responses "echo-api/models/dtos/responses/pagination"
	"echo-api/models/entities"
	"echo-api/util"
	"fmt"

	"gorm.io/gorm/clause"
)

type LanguageService struct {
	repo   util.Repository[entities.Language]
	logger *util.Logger
}

func NewLanguageService(repo util.Repository[entities.Language], logger *util.Logger) *LanguageService {
	return &LanguageService{repo: repo, logger: logger}
}

func (s *LanguageService) GetOne(id string) (entities.Language, error) {
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_GetOne with id: %s", id))
	language, err := s.repo.First(id, true)
	if err != nil {
		s.logger.Error().Msg("LanguageService_GetOne had an error when saving to repo")
		return entities.Language{}, err
	}

	return language, nil
}

func (s *LanguageService) FilterAll(request requests.FilterLanguagesRequest) (responses.PaginationResponse[entities.Language], error) {
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	languages, err := q.Find(true)
	if err != nil {
		s.logger.Error().Msg("LanguageService_FilterAll had an error when requesting the data from repo")
		return responses.PaginationResponse[entities.Language]{}, err
	}
	count, err := q.Count()
	if err != nil {
		s.logger.Error().Msg("LanguageService_FilterAll had an error when requesting the data from repo")
		return responses.PaginationResponse[entities.Language]{}, err
	}
	return responses.PaginationResponse[entities.Language]{Content: languages, Page: request.Page, Size: len(languages), TotalCount: int(count)}, nil
}

func (s *LanguageService) buildFilterQuery(request requests.FilterLanguagesRequest) util.Repository[entities.Language] {
	q := s.repo.Query()
	s.logger.Debug().Msg("*LanguageService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("*LanguageService filtering IDs")
		q = q.Where("ID IN ?", *request.IDs)
	}

	if request.Name != nil && len(*request.Name) > 0 {
		s.logger.Debug().Msg("*LanguageService filtering Name")
		nameLike := "%" + *request.Name + "%"
		q = q.Where("name LIKE ?", nameLike)
	}

	if request.UserIDs != nil && len(*request.UserIDs) > 0 {
		s.logger.Debug().Msg("*LanguageService filtering languages")
		q = q.Where("userID IN ?", *request.UserIDs)
	}

	if request.CourseIDs != nil && len(*request.CourseIDs) > 0 {
		s.logger.Debug().Msg("*LanguageService filtering Courses")
		q = q.Where("courseID IN ?", *request.CourseIDs)
	}

	if request.Alpha2Code != nil && len(*request.Alpha2Code) > 0 {
		s.logger.Debug().Msg("*LanguageService filtering languages")
		q.Where("alpha2Code = ?", *request.Alpha2Code)
	}

	return q.Order("created_at")
}

func (s *LanguageService) CreateOne(request requests.CreateLanguageRequest) (entities.Language, error) {
	s.logger.Debug().Msg("LanguageService_CreateOne has started")
	language := entities.Language{
		Name:       request.Name,
		Alpha2Code: request.Alpha2Code,
		Alpha3Code: request.Alpha3Code,
		Icon:       request.Icon,
	}
	language, err := s.repo.Create(&language)
	if err != nil {
		s.logger.Error().Msg("LanguageService_CreateOne had an error when saving to repo")
		return entities.Language{}, err
	}
	return language, nil
}

func (s *LanguageService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_DeleteOne has started with given id: %s", id))
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error().Msg("LanguageService_DeleteOne had an error when deleting from repo")
		return false, err
	}

	return true, nil
}

func (s *LanguageService) UpdateOne(request requests.UpdateLanguageRequest) (entities.Language, error) {
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_UpdateOne has started with given id: %s", request.ID))
	language, err := s.repo.Clauses(clause.Locking{Strength: "UPDATE"}).First(request.ID, true)
	if err != nil {
		s.logger.Error().Msg(fmt.Sprintf("LanguageService_UpdateOne could not find a record with given id: %s", request.ID))
		return entities.Language{}, err
	}

	if request.Name != nil && *request.Name != "" {
		s.logger.Debug().Msg(fmt.Sprintf("LanguageService_UpdateOne updating Name. From: %v => To: %v", language.Name, *request.Name))
		language.Name = *request.Name
	}

	if request.Alpha2Code != nil && *request.Alpha2Code != "" {
		s.logger.Debug().Msg(fmt.Sprintf("LanguageService_UpdateOne updating Alpha2Code. From: %v => To: %v", language.Alpha2Code, *request.Alpha2Code))
		language.Alpha2Code = *request.Alpha2Code
	}

	if request.Alpha3Code != nil && *request.Alpha3Code != "" {
		s.logger.Debug().Msg(fmt.Sprintf("LanguageService_UpdateOne updating Alpha3Code. From: %v => To: %v", language.Alpha3Code, *request.Alpha3Code))
		language.Alpha3Code = *request.Alpha3Code
	}

	language, err = s.repo.Update(&language)
	if err != nil {
		s.logger.Error().Msg("LanguageService_UpdateOne had an error while trying to save to repo")
		return entities.Language{}, err
	}
	return language, nil
}
