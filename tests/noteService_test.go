package tests

import (
	"echo-api/mocks"
	"echo-api/models/dtos/requests/base"
	"echo-api/models/dtos/requests/note"
	"echo-api/models/entities"
	"echo-api/services"
	"echo-api/util"
	"os"
	"testing"
)

func TestCreateNote_Success(t *testing.T) {
	s := getMockedNoteService()
	one := "1"
	req := note.CreateNoteRequest{
		Header:     "Test",
		Payload:    "test1",
		LanguageID: "1",
		UserID:     &one,
		ContextID:  "1",
	}

	note, err := s.CreateOne(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.Header != note.Header {
		t.Errorf("Expected %s but got %s", req.Header, note.Header)
		return
	}
	if req.Payload != note.Payload {
		t.Errorf("Expected %s but got %s", req.Payload, note.Payload)
		return
	}
}

func TestGetOneNoteWithID_Success(t *testing.T) {
	s := getMockedNoteService()
	one := "1"
	req := note.CreateNoteRequest{
		Header:     "Test",
		Payload:    "test1",
		LanguageID: "1",
		UserID:     &one,
		ContextID:  "1",
	}

	note, err := s.CreateOne(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	result, err := s.GetOne(note.ID)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.Header != result.Header {
		t.Errorf("Expected %s but got %s", req.Header, result.Header)
		return
	}
	if req.Payload != result.Payload {
		t.Errorf("Expected %s but got %s", req.Payload, result.Payload)
		return
	}
}

func TestFilterAllNotes_Success(t *testing.T) {
	s := getMockedNoteService()
	one := "1"
	notes := []note.CreateNoteRequest{
		{
			Header:     "Test",
			Payload:    "test1",
			LanguageID: "1",
			UserID:     &one,
			ContextID:  "1",
		},
		{
			Header:     "Test Test",
			Payload:    "test2",
			LanguageID: "1",
			UserID:     &one,
			ContextID:  "2",
		},
		{
			Header:     "Test Test Test",
			Payload:    "test3",
			LanguageID: "2",
			UserID:     &one,
			ContextID:  "3",
		},
	}

	for _, v := range notes {
		_, err := s.CreateOne(v)
		if err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
			return
		}
	}

	request := note.FilterNotesRequest{PaginationRequestBase: base.PaginationRequestBase{Page: 0, Size: 10}}

	res, err := s.FilterAll(request)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if res.Size != len(notes) {
		t.Errorf("Result len is not matching the expected len: %d", len(notes))
		return
	}

	if res.Page != request.Page {
		t.Errorf("Expected %d but got %d", request.Page, res.Page)
		return
	}
	// TODO check all exists at least once
	//
	//	if res.Content[0].Email != notes[0].Email {
	//		t.Errorf("Expected %s but got %s", notes[0].Email, res.Content[0].Email)
	//		return
	//	}
}

func TestDeleteNote_Success(t *testing.T) {
	s := getMockedNoteService()
	one := "1"
	notes := []note.CreateNoteRequest{
		{
			Header:     "Test",
			Payload:    "test1",
			LanguageID: "1",
			UserID:     &one,
			ContextID:  "1",
		},
		{
			Header:     "Test Test",
			Payload:    "test2",
			LanguageID: "1",
			UserID:     &one,
			ContextID:  "2",
		},
	}
	for _, v := range notes {
		_, err := s.CreateOne(v)
		if err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
			return
		}
	}

	// Mock db always give 1 to first Id so we expect first elem to get deleted
	_, err := s.DeleteOne("1")
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	res, err := s.FilterAll(note.FilterNotesRequest{PaginationRequestBase: base.PaginationRequestBase{Page: 1, Size: 100}})
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if res.Size > 1 || res.Size == 0 {
		t.Errorf("Expected a single element but got %d", res.Size)
		return
	}

	if res.Content[0].Payload != notes[1].Payload {
		t.Errorf("Expected %s but got %s", notes[1].Payload, res.Content[0].Payload)
		return
	}
}

func TestDeleteNote_NotFoundError(t *testing.T) {
	s := getMockedNoteService()
	one := "1"
	notes := []note.CreateNoteRequest{
		{
			Header:     "Test",
			Payload:    "test1",
			LanguageID: "1",
			UserID:     &one,
			ContextID:  "1",
		},
		{
			Header:     "Test Test",
			Payload:    "test2",
			LanguageID: "1",
			UserID:     &one,
			ContextID:  "2",
		},
	}
	for _, v := range notes {
		_, err := s.CreateOne(v)
		if err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
			return
		}
	}

	_, err := s.DeleteOne("3")
	if err == nil {
		t.Errorf("Expected errors but got none")
		return
	}

	if err.Error() != "notFoundError" {
		t.Errorf("Expected \"notFoundError\" but got %s", err.Error())
		return
	}
}

func getMockedNoteService() *services.NoteService {
	mockRepo := mocks.NewMockRepo[entities.Note]()
	logger := util.NewLogger(map[string]string{}, os.Stdout)
	return services.NewNoteService(mockRepo, logger)
}
