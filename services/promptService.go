package services

import (
	"echo-api/managers"
	requests "echo-api/models/dtos/requests/prompt"
	responses "echo-api/models/dtos/responses/pagination"
	"echo-api/models/entities"
	"echo-api/util"
	"errors"
	"fmt"
)

type PromptService struct {
	repo          util.Repository[entities.Prompt]
	logger        *util.Logger
	promptManager managers.PromptGenManager
	commsManager  managers.AiCommunicationManager
}

func NewPromptService(repo util.Repository[entities.Prompt], logger *util.Logger, pm managers.PromptGenManager, cm managers.AiCommunicationManager) *PromptService {
	return &PromptService{repo: repo, logger: logger, promptManager: pm, commsManager: cm}
}

func (s *PromptService) GetOne(id string) (entities.Prompt, error) {
	if id == "" {
		return entities.Prompt{}, errors.New("argumentErrorIDMissing")
	}
	s.logger.Debug().Msg(fmt.Sprintf("PromptService_GetOne with id: %s", id))
	res, err := s.repo.First(id, false)
	if err != nil {
		s.logger.Error().Msg("PromptService_GetOne had an error when saving to repo")
		return entities.Prompt{}, err
	}

	return res, nil
}

func (s *PromptService) FilterAll(request requests.FilterPromptsRequest) (responses.PaginationResponse[entities.Prompt], error) {
	s.logger.Debug().Msg(fmt.Sprintf("PromptService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	res, err := q.Find(false)
	if err != nil {
		s.logger.Error().Msg("PromptService_FilterAll had an error when requesting from repo")
		return responses.PaginationResponse[entities.Prompt]{}, err
	}
	count, err := q.Count()
	if err != nil {
		s.logger.Error().Msg("PromptService_FilterAll had an error when requesting from repo")
		return responses.PaginationResponse[entities.Prompt]{}, err
	}
	return responses.PaginationResponse[entities.Prompt]{Content: res, Page: request.Page, Size: len(res), TotalCount: int(count)}, nil
}

func (s *PromptService) buildFilterQuery(request requests.FilterPromptsRequest) util.Repository[entities.Prompt] {
	q := s.repo.Query()
	s.logger.Debug().Msg("*PromptService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("*PromptService filtering IDs")
		q = q.Where("ID IN ?", *request.IDs)
	}

	if request.ContextIDs != nil && len(*request.ContextIDs) > 0 {
		s.logger.Debug().Msg("*PromptService filtering ContextID")
		q = q.Where("contextID IN ?", *request.ContextIDs)
	}

	return q.Order("created_at")
}

func (s *PromptService) FindByEntityAndContext(request requests.FindPromptByEntityAndContextRequest) (entities.Prompt, error) {
	if request.EntityID == "" {
		return entities.Prompt{}, errors.New("argumentErrorIDMissing")
	}
	q := s.repo.Query()
	if request.ContextID != "" {
		q = q.Where("contextID = ?", request.ContextID)
	}
	res, err := q.Where("entityID = ?", request.EntityID).Find(false)
	if err != nil {
		s.logger.Error().Msg("PromptService_GetOne had an error when fetching from repo")
		return entities.Prompt{}, err
	} else if len(res) == 0 {
		return entities.Prompt{}, errors.New("notFoundError")
	}

	return res[0], nil
}

func (s *PromptService) GenerateAndSendPrompt(request requests.CreatePromptRequest) (entities.Prompt, error) {
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
	prompt, err = s.repo.Create(&prompt)
	if err != nil {
		s.logger.Error().Msg("PromptService_CreateOne had an error when saving to repo")
		return entities.Prompt{}, err
	}

	return prompt, nil
}

func (s *PromptService) GenerateAndSendMessage(request requests.CreateMessageRequest) (string, error) {
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

func (s *PromptService) DeleteAndSend(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("PromptService_DeleteOne has started with given id: %s", id))
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error().Msg("PromptService_DeleteOne had an error when deleting from repo")
		return false, err
	}

	return true, nil
}

func (s *PromptService) UpdatePrompt(request requests.UpdatePromptRequest) (entities.Prompt, error) {
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
