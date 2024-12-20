package services

import (
	"errors"
	"fmt"
	requests "reson8-learning-api/models/dtos/requests/note"
	responses "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NoteService struct {
	db     *gorm.DB
	logger *util.Logger
}

func NewNoteService(db *gorm.DB, logger *util.Logger) *NoteService {
	return &NoteService{db: db, logger: logger}
}

func (s *NoteService) CheckIfBelongsToUser(id string, userID string) (bool, error) {
	var dbID string
	s.logger.Debug().Msg(fmt.Sprintf("NoteService_CheckIfBelongsToUser with id: %s for user: %s", id, userID))
	res := s.db.Model(entities.Note{}).Pluck("userID", &dbID)
	if res.Error != nil {
		s.logger.Error().Msg("NoteService_CheckIfBelongsToUser had an error when getting from db")
		return false, res.Error
	}

	return userID == dbID, nil
}

func (s NoteService) GetOne(id string) (entities.Note, error) {
	var note entities.Note
	s.logger.Debug().Msg(fmt.Sprintf("NoteService_GetOne with id: %s", id))

	res := s.db.First(&note, id)
	if res.Error != nil {
		s.logger.Error().Msg("NoteService_GetOne had an error when saving to db")
		return entities.Note{}, res.Error
	}

	return note, nil
}

func (s NoteService) FilterAll(request requests.FilterNotesRequest) (responses.PaginationResponse[entities.Note], error) {
	var notes []entities.Note
	s.logger.Debug().Msg(fmt.Sprintf("NoteService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	res := q.Find(&notes)
	if res.Error != nil {
		s.logger.Error().Msg("NoteService_FilterAll had an error when requesting from db")
		return responses.PaginationResponse[entities.Note]{}, res.Error
	}
	var count int64
	res = q.Count(&count)
	if res.Error != nil {
		s.logger.Error().Msg("NoteService_FilterAll had an error when requesting from db")
		return responses.PaginationResponse[entities.Note]{}, res.Error
	}
	return responses.PaginationResponse[entities.Note]{Content: notes, Page: request.Page, Size: len(notes), TotalCount: int(count)}, nil
}

func (s NoteService) buildFilterQuery(request requests.FilterNotesRequest) *gorm.DB {
	q := s.db.Model(&entities.Note{})
	s.logger.Debug().Msg("NoteService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("NoteService filtering IDs")
		q = q.Where(*request.IDs)
	}

	if request.Header != nil && *request.Header != "" {
		s.logger.Debug().Msg("NoteService filtering Header")
		headerLike := "%" + *request.Header + "%"
		q = q.Where("header LIKE ?", headerLike)
	}

	if request.UserIDs != nil && len(*request.UserIDs) > 0 {
		s.logger.Debug().Msg("NoteService filtering UserIDs")
		q = q.Where("userID IN ?", *request.UserIDs)
	}

	if request.LanguageIDs != nil && len(*request.LanguageIDs) > 0 {
		s.logger.Debug().Msg("NoteService filtering LanguageIDs")
		q = q.Where("languageID IN ?", *request.LanguageIDs)
	}

	if request.DocumentIDs != nil && len(*request.DocumentIDs) > 0 {
		s.logger.Debug().Msg("NoteService filtering DocumentID")
		q = q.Where("documentID IN ?", *request.DocumentIDs)
	}

	if request.ContextIDs != nil && len(*request.ContextIDs) > 0 {
		s.logger.Debug().Msg("NoteService filtering ContextID")
		q = q.Where("contextID IN ?", *request.ContextIDs)
	}

	return q.Order("created_at")
}

func (s NoteService) CreateOne(request requests.CreateNoteRequest) (entities.Note, error) {
	if request.UserID == nil || *request.UserID == "" {
		return entities.Note{}, errors.New("argumentErrorIDMissing")
	}
	note := entities.Note{
		Header:     request.Header,
		Payload:    request.Payload,
		UserID:     *request.UserID,
		LanguageID: request.LanguageID,
		ContextID:  request.ContextID,
	}
	s.logger.Debug().Msg("NoteService_CreateOne has started")
	res := s.db.Create(&note)
	if res.Error != nil {
		s.logger.Error().Msg("NoteService_CreateOne had an error when saving to db")
		return entities.Note{}, res.Error
	}

	return note, nil
}

func (s NoteService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("NoteService_DeleteOne has started with given id: %s", id))
	res := s.db.Select(clause.Associations).Delete(&entities.Note{}, id)
	if res.Error != nil {
		s.logger.Error().Msg("NoteService_DeleteOne had an error when deleting from db")
		return false, res.Error
	}

	return true, nil
}

func (s NoteService) UpdateOne(request requests.UpdateNoteRequest) (entities.Note, error) {
	s.logger.Debug().Msg(fmt.Sprintf("NoteService_UpdateOne has started with given id: %s", request.ID))
	var note entities.Note
	res := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&note, request.ID)
	err := res.Error
	if err != nil {
		s.logger.Error().Msg(fmt.Sprintf("NoteService_UpdateOne could not find a record with given id: %s", request.ID))
		return entities.Note{}, err
	}

	if request.Header != nil && *request.Header != "" {
		s.logger.Debug().Msg(fmt.Sprintf("NoteService_UpdateOne updated Header. From: %v => To: %v", note.Header, *request.Header))
		note.Header = *request.Header
	}

	if request.Payload != nil && *request.Payload != "" {
		s.logger.Debug().Msg(fmt.Sprintf("NoteService_UpdateOne updated Payload. From: %v => To: %v", note.Payload, *request.Payload))
		note.Payload = *request.Payload
	}

	if request.LanguageID != nil && *request.LanguageID != note.LanguageID {
		s.logger.Debug().Msg(fmt.Sprintf("NoteService_UpdateOne updated Language. From: %v => To: %v", note.LanguageID, *request.LanguageID))
		note.LanguageID = *request.LanguageID
	}

	if request.UserID != nil && *request.UserID != note.UserID {
		s.logger.Debug().Msg(fmt.Sprintf("NoteService_UpdateOne updated User. From: %v => To: %v", note.UserID, *request.UserID))
		note.UserID = *request.UserID
	}

	res = s.db.Save(&note)
	if res.Error != nil {
		s.logger.Error().Msg("NoteService_UpdateOne had an error while trying to save to db")
		return entities.Note{}, res.Error
	}
	return note, nil
}
