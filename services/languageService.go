package services

import (
	"fmt"
	requests "reson8-learning-api/models/dtos/requests/language"
	responses "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LanguageService struct {
	db     *gorm.DB
	logger *util.Logger
}

func NewLanguageService(db *gorm.DB, logger *util.Logger) *LanguageService {
	return &LanguageService{db: db, logger: logger}
}

func (s LanguageService) GetOne(id string) (entities.Language, error) {
	var language entities.Language
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_GetOne with id: %s", id))
	res := s.db.First(&language, id)
	if res.Error != nil {
		s.logger.Error().Msg("LanguageService_GetOne had an error when saving to db")
		return entities.Language{}, res.Error
	}

	return language, nil
}

func (s LanguageService) FilterAll(request requests.FilterLanguagesRequest) (responses.PaginationResponse[entities.Language], error) {
	var languages []entities.Language
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	res := q.Find(&languages)
	if res.Error != nil {
		s.logger.Error().Msg("LanguageService_FilterAll had an error when requesting the data from db")
		return responses.PaginationResponse[entities.Language]{}, res.Error
	}
	var count int64
	res = q.Count(&count)
	if res.Error != nil {
		s.logger.Error().Msg("LanguageService_FilterAll had an error when requesting the data from db")
		return responses.PaginationResponse[entities.Language]{}, res.Error
	}
	return responses.PaginationResponse[entities.Language]{Content: languages, Page: request.Page, Size: len(languages), TotalCount: int(count)}, nil
}

func (s LanguageService) buildFilterQuery(request requests.FilterLanguagesRequest) *gorm.DB {
	q := s.db.Model(&entities.Language{})
	s.logger.Debug().Msg("LanguageService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("LanguageService filtering IDs")
		q = q.Where(*request.IDs)
	}

	if request.Name != nil && len(*request.Name) > 0 {
		s.logger.Debug().Msg("LanguageService filtering Name")
		nameLike := "%" + *request.Name + "%"
		q = q.Where("name LIKE ?", nameLike)
	}

	if request.UserIDs != nil && len(*request.UserIDs) > 0 {
		s.logger.Debug().Msg("LanguageService filtering languages")
		q = q.Preload("Users")
		q = q.Where("userID IN ?", *request.UserIDs)
	}

	if request.CourseIDs != nil && len(*request.CourseIDs) > 0 {
		s.logger.Debug().Msg("LanguageService filtering Courses")
		q = q.Preload("Courses")
		q = q.Where("courseID IN ?", *request.CourseIDs)
	}

	if request.Alpha2Code != nil && len(*request.Alpha2Code) > 0 {
		s.logger.Debug().Msg("LanguageService filtering languages")
		q.Where("alpha2Code = ?", *request.Alpha2Code)
	}

	return q.Order("created_at")
}

func (s LanguageService) CreateOne(request requests.CreateLanguageRequest) (entities.Language, error) {
	language := entities.Language{
		Name:       request.Name,
		Alpha2Code: request.Alpha2Code,
		Alpha3Code: request.Alpha3Code,
		Icon:       request.Icon,
	}
	s.logger.Debug().Msg("LanguageService_CreateOne has started")
	res := s.db.Create(&language)
	if res.Error != nil {
		s.logger.Error().Msg("LanguageService_CreateOne had an error when saving to db")
		return entities.Language{}, res.Error
	}
	return language, nil
}

func (s LanguageService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_DeleteOne has started with given id: %s", id))
	res := s.db.Delete(&entities.Language{}, id)
	if res.Error != nil {
		s.logger.Error().Msg("LanguageService_DeleteOne had an error when deleting from db")
		return false, res.Error
	}

	return true, nil
}

func (s LanguageService) UpdateOne(request requests.UpdateLanguageRequest) (entities.Language, error) {
	s.logger.Debug().Msg(fmt.Sprintf("LanguageService_UpdateOne has started with given id: %s", request.ID))
	var language entities.Language
	res := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&language, request.ID)
	err := res.Error
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

	res = s.db.Save(&language)
	if res.Error != nil {
		s.logger.Error().Msg("LanguageService_UpdateOne had an error while trying to save to db")
		return entities.Language{}, res.Error
	}
	return language, nil
}
