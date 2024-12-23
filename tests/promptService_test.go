package tests

import (
	"echo-api/managers/implementations"
	"echo-api/mocks"
	"echo-api/models/dtos/requests/base"
	"echo-api/models/dtos/requests/prompt"
	"echo-api/models/entities"
	"echo-api/services"
	"echo-api/util"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCreatePromptFromString_Success(t *testing.T) {
	req := prompt.CreatePromptRequest{
		Value:     "XXX",
		ContextID: "1",
		EntityID:  "1",
	}

	s := getMockedPromptService()
	prompt, err := s.GenerateAndSendPrompt(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.ContextID != prompt.ContextID {
		t.Errorf("Expected %s but got %s", req.ContextID, prompt.ContextID)
		return
	}
	if req.EntityID != *prompt.EntityID {
		t.Errorf("Expected %s but got %s", req.EntityID, *prompt.EntityID)
		return
	}
	if prompt.Value != "Prompt(Remember(XXX))" {
		t.Errorf("Expected %s but got %s", "Prompt(Remember(XXX))", prompt.Value)
		return
	}
}

func TestCreatePromptFromNote_Success(t *testing.T) {
	req := prompt.CreatePromptRequest{
		Value: entities.Note{
			Header:     "Test",
			Payload:    "Testy Test",
			UserID:     "1",
			LanguageID: "1",
			ContextID:  "1",
			Base: entities.Base{
				ID:        "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		ContextID: "1",
		EntityID:  "1",
	}

	s := getMockedPromptService()
	prompt, err := s.GenerateAndSendPrompt(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.ContextID != prompt.ContextID {
		t.Errorf("Expected %s but got %s", req.ContextID, prompt.ContextID)
		return
	}
	if req.EntityID != *prompt.EntityID {
		t.Errorf("Expected %s but got %s", req.EntityID, *prompt.EntityID)
		return
	}
	expectedVal := fmt.Sprintf("Prompt(Remember(%s))", req.Value.(entities.Note).Payload)
	if prompt.Value != expectedVal {
		t.Errorf("Expected %s but got %s", expectedVal, prompt.Value)
		return
	}
}

func TestCreatePromptFromDocument_Success(t *testing.T) {
	t.Skip("Skipping this test until implementing mock for file manager")
	req := prompt.CreatePromptRequest{
		Value: entities.Document{
			Name:      "createPromptTest_1.pdf",
			Location:  "Pdf",
			UserID:    "1",
			Extension: "pdf",
			ContextID: "1",
			Base: entities.Base{
				ID:        "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		ContextID: "1",
		EntityID:  "1",
	}

	s := getMockedPromptService()
	prompt, err := s.GenerateAndSendPrompt(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.ContextID != prompt.ContextID {
		t.Errorf("Expected %s but got %s", req.ContextID, prompt.ContextID)
		return
	}
	if req.EntityID != *prompt.EntityID {
		t.Errorf("Expected %s but got %s", req.EntityID, *prompt.EntityID)
		return
	}
	expectedVal := fmt.Sprintf("Prompt(Remember(%s))", req.Value.(entities.Note).Payload)
	if prompt.Value != expectedVal {
		t.Errorf("Expected %s but got %s", expectedVal, prompt.Value)
		return
	}
}

func TestGetOnePrompt_Success(t *testing.T) {
	req := prompt.CreatePromptRequest{
		Value:     "XXX",
		ContextID: "1",
		EntityID:  "1",
	}
	s := getMockedPromptService()
	prompt, err := s.GenerateAndSendPrompt(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	result, err := s.GetOne(prompt.ID)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.ContextID != result.ContextID {
		t.Errorf("Expected %s but got %s", req.ContextID, result.ContextID)
		return
	}
	if req.EntityID != *result.EntityID {
		t.Errorf("Expected %s but got %s", req.EntityID, *result.EntityID)
		return
	}
	if prompt.Value != "Prompt(Remember(XXX))" {
		t.Errorf("Expected %s but got %s", "Prompt(Remember(XXX))", result.Value)
		return
	}
}

func TestFilterAllPrompts_Success(t *testing.T) {
	s := getMockedPromptService()
	prompts := []prompt.CreatePromptRequest{
		{
			Value:     "XXX",
			ContextID: "1",
			EntityID:  "2",
		},
		{
			Value:     "ZZZ",
			ContextID: "1",
			EntityID:  "3",
		},
		{
			Value:     "TTT",
			ContextID: "2",
			EntityID:  "10",
		},
	}

	for _, v := range prompts {
		_, err := s.GenerateAndSendPrompt(v)
		if err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
			return
		}
	}

	request := prompt.FilterPromptsRequest{PaginationRequestBase: base.PaginationRequestBase{Page: 0, Size: 10}}

	res, err := s.FilterAll(request)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if res.Size != len(prompts) {
		t.Errorf("Result len is not matching the expected len: %d", len(prompts))
		return
	}

	if res.Page != request.Page {
		t.Errorf("Expected %d but got %d", request.Page, res.Page)
		return
	}
}

func TestDeletePrompt_Success(t *testing.T) {
	s := getMockedPromptService()
	prompts := []prompt.CreatePromptRequest{
		{
			Value:     "XXX",
			ContextID: "1",
			EntityID:  "2",
		},
		{
			Value:     "ZZZ",
			ContextID: "1",
			EntityID:  "3",
		},
	}
	for _, v := range prompts {
		_, err := s.GenerateAndSendPrompt(v)
		if err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
			return
		}
	}

	// Mock db always give 1 to first Id so we expect first elem to get deleted
	_, err := s.DeleteAndSend("1")
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	res, err := s.FilterAll(prompt.FilterPromptsRequest{PaginationRequestBase: base.PaginationRequestBase{Page: 1, Size: 100}})
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if res.Size > 1 || res.Size == 0 {
		t.Errorf("Expected a single element but got %d", res.Size)
		return
	}

	if res.Content[0].Value != "Prompt(Remember(ZZZ))" {
		t.Errorf("Expected %s but got %s", "Prompt(Remember(ZZZ))", res.Content[0].Value)
		return
	}
}

func TestDeletePrompt_NotFoundError(t *testing.T) {
	s := getMockedPromptService()
	prompts := []prompt.CreatePromptRequest{
		{
			Value:     "XXX",
			ContextID: "1",
			EntityID:  "2",
		},
		{
			Value:     "ZZZ",
			ContextID: "1",
			EntityID:  "3",
		},
	}
	for _, v := range prompts {
		_, err := s.GenerateAndSendPrompt(v)
		if err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
			return
		}
	}

	_, err := s.DeleteAndSend("3")
	if err == nil {
		t.Errorf("Expected errors but got none")
		return
	}

	if err.Error() != "notFoundError" {
		t.Errorf("Expected \"notFoundError\" but got %s", err.Error())
		return
	}
}

func getMockedPromptService() *services.PromptService {
	mockRepo := mocks.NewMockRepo[entities.Prompt]()
	logger := util.NewLogger(map[string]string{}, os.Stdout)
	fileManager := implementations.NewOnServerFileManager("assets", []string{})
	promptGenManager := implementations.NewLocalPromptGenManager(fileManager)
	aiCommManager := implementations.NewOpenAiCommunicationManager(&util.Configuration{IsAiAssistantEnabled: false})
	return services.NewPromptService(mockRepo, logger, promptGenManager, aiCommManager)
}
