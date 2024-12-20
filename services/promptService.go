package services

import (
	"errors"
	"fmt"
	"reson8-learning-api/managers"
	requests "reson8-learning-api/models/dtos/requests/prompt"
	responses "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PromptService struct {
	db            *gorm.DB
	logger        *util.Logger
	promptManager managers.PromptGenManager
	commsManager  managers.AiCommunicationManager
}

func NewPromptService(db *gorm.DB, logger *util.Logger, pm managers.PromptGenManager, cm managers.AiCommunicationManager) *PromptService {
	return &PromptService{db: db, logger: logger, promptManager: pm, commsManager: cm}
}

func (s PromptService) GetOne(id string) (entities.Prompt, error) {
	if id == "" {
		return entities.Prompt{}, errors.New("argumentErrorIDMissing")
	}

	var prompt entities.Prompt
	s.logger.Debug().Msg(fmt.Sprintf("PromptService_GetOne with id: %s", id))

	res := s.db.First(&prompt, id)
	if res.Error != nil {
		s.logger.Error().Msg("PromptService_GetOne had an error when saving to db")
		return entities.Prompt{}, res.Error
	}

	return prompt, nil
}

func (s PromptService) FilterAll(request requests.FilterPromptsRequest) (responses.PaginationResponse[entities.Prompt], error) {
	var prompts []entities.Prompt
	s.logger.Debug().Msg(fmt.Sprintf("PromptService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	res := q.Find(&prompts)
	if res.Error != nil {
		s.logger.Error().Msg("PromptService_FilterAll had an error when requesting from db")
		return responses.PaginationResponse[entities.Prompt]{}, res.Error
	}
	var count int64
	res = q.Count(&count)
	if res.Error != nil {
		s.logger.Error().Msg("PromptService_FilterAll had an error when requesting from db")
		return responses.PaginationResponse[entities.Prompt]{}, res.Error
	}
	return responses.PaginationResponse[entities.Prompt]{Content: prompts, Page: request.Page, Size: len(prompts), TotalCount: int(count)}, nil
}

func (s PromptService) buildFilterQuery(request requests.FilterPromptsRequest) *gorm.DB {
	q := s.db.Model(&entities.Prompt{})
	s.logger.Debug().Msg("PromptService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("PromptService filtering IDs")
		q = q.Where(*request.IDs)
	}

	if request.ContextIDs != nil && len(*request.ContextIDs) > 0 {
		s.logger.Debug().Msg("PromptService filtering ContextID")
		q = q.Where("contextID IN ?", *request.ContextIDs)
	}

	return q.Order("created_at")
}

func (s PromptService) FindByEntityAndContext(request requests.FindPromptByEntityAndContextRequest) (entities.Prompt, error) {
	if request.EntityID == "" {
		return entities.Prompt{}, errors.New("argumentErrorIDMissing")
	}
	var prompt entities.Prompt

	q := s.db
	if request.ContextID != "" {
		q = q.Where("contextID = ?")
	}
	res := q.First(&prompt, "entityID = ?", request.ContextID, request.EntityID)
	if res.Error != nil {
		s.logger.Error().Msg("PromptService_GetOne had an error when fetching from db")
		return entities.Prompt{}, res.Error
	}

	return prompt, nil
}

func (s PromptService) GenerateAndSendPrompt(request requests.CreatePromptRequest) (entities.Prompt, error) {
	if request.ContextID == "" {
		return entities.Prompt{}, errors.New("argumentErrorIDMissing")
	}
	promptValue, err := s.promptManager.GeneratePrompt(request.Value)
	if err != nil {
		return entities.Prompt{}, err
	}
	_, err = s.commsManager.SendPrompt(request.ContextID, promptValue)
	if err != nil {
		return entities.Prompt{}, err
	}

	prompt := entities.Prompt{
		ContextID: request.ContextID,
		Value:     promptValue,
		EntityID:  &request.EntityID,
	}
	if request.EntityID == "" {
		prompt.EntityID = nil
	}
	s.logger.Debug().Msg("PromptService_CreateOne has started")
	res := s.db.Create(&prompt)
	if res.Error != nil {
		s.logger.Error().Msg("PromptService_CreateOne had an error when saving to db")
		return entities.Prompt{}, res.Error
	}

	return prompt, nil
}

func (s PromptService) GenerateAndSendMessage(request requests.CreateMessageRequest) (string, error) {
	if request.ContextID == "" {
		return "", errors.New("argumentErrorIDMissing")
	}
	promptValue, err := s.promptManager.GenerateMessage(request.Value)
	if err != nil {
		return "", err
	}
	resp, err := s.commsManager.SendPrompt(request.ContextID, promptValue)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (s PromptService) DeleteAndSend(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("PromptService_DeleteOne has started with given id: %s", id))
	res := s.db.Select(clause.Associations).Delete(&entities.Prompt{}, id)
	if res.Error != nil {
		s.logger.Error().Msg("PromptService_DeleteOne had an error when deleting from db")
		return false, res.Error
	}

	return true, nil
}

func (s PromptService) UpdatePrompt(request requests.UpdatePromptRequest) (entities.Prompt, error) {
	findRequest := requests.FindPromptByEntityAndContextRequest{
		EntityID:  request.EntityID,
		ContextID: request.ContextID,
	}
	found, err := s.FindByEntityAndContext(findRequest)
	if err != nil {
		return entities.Prompt{}, err
	}
	_, err = s.DeleteAndSend(found.ID)
	if err != nil {
		return entities.Prompt{}, err
	}

	return s.GenerateAndSendPrompt(requests.CreatePromptRequest(request))
}
