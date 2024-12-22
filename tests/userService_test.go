package tests

import (
	"os"
	"reson8-learning-api/mocks"
	"reson8-learning-api/models/dtos/requests/base"
	"reson8-learning-api/models/dtos/requests/user"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/services"
	"reson8-learning-api/util"
	"testing"
)

func TestCreateUserWithCorrectArgs(t *testing.T) {
	req := user.CreateUserRequest{
		Name:     "XXX YYY",
		Email:    "example@mail.com",
		Password: "!testPass_4251",
	}

	s := getMockedUserService()
	user, err := s.CreateOne(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.Email != user.Email {
		t.Errorf("Expected %s but got %s", req.Email, user.Email)
		return
	}
	if req.Name != user.Name {
		t.Errorf("Expected %s but got %s", req.Name, user.Name)
		return
	}
}

func TestGetOneUserWithID(t *testing.T) {
	req := user.CreateUserRequest{
		Name:     "XXX YYY",
		Email:    "example@mail.com",
		Password: "!testPass_4251",
	}
	s := getMockedUserService()
	user, err := s.CreateOne(req)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	result, err := s.GetOne(user.ID)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if req.Email != result.Email {
		t.Errorf("Expected %s but got %s", req.Email, result.Email)
		return
	}
	if req.Name != result.Name {
		t.Errorf("Expected %s but got %s", req.Name, result.Name)
		return
	}
}

func TestFilterAll(t *testing.T) {
	s := getMockedUserService()
	users := []user.CreateUserRequest{
		{
			Name:     "XXX YYY",
			Email:    "example1@mail.com",
			Password: "!testPass_4251",
		},
		{
			Name:     "XXX ZZZ",
			Email:    "example2@mail.com",
			Password: "!testPass_4251",
		},
		{
			Name:     "XXX TTT",
			Email:    "example3@mail.com",
			Password: "!testPass_4251",
		},
	}

	for _, v := range users {
		_, err := s.CreateOne(v)
		if err != nil {
			t.Errorf("Expected no errors but got %s", err.Error())
			return
		}
	}

	request := user.FilterUsersRequest{PaginationRequestBase: base.PaginationRequestBase{Page: 0, Size: 10}}

	res, err := s.FilterAll(request)
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if res.Size != len(users) {
		t.Errorf("Result len is not matching the expected len: %d", len(users))
		return
	}

	if res.Page != request.Page {
		t.Errorf("Expected %d but got %d", request.Page, res.Page)
		return
	}

	if res.Content[0].Email != users[0].Email {
		t.Errorf("Expected %s but got %s", users[0].Email, res.Content[0].Email)
		return
	}
}

func TestDeleteOne(t *testing.T) {
	s := getMockedUserService()
	users := []user.CreateUserRequest{
		{
			Name:     "XXX YYY",
			Email:    "example1@mail.com",
			Password: "!testPass_4251",
		},
		{
			Name:     "XXX ZZZ",
			Email:    "example2@mail.com",
			Password: "!testPass_4251",
		},
	}
	for _, v := range users {
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

	res, err := s.FilterAll(user.FilterUsersRequest{PaginationRequestBase: base.PaginationRequestBase{Page: 1, Size: 100}})
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
		return
	}

	if res.Size > 1 || res.Size == 0 {
		t.Errorf("Expected a single element but got %d", res.Size)
		return
	}

	if res.Content[0].Email != users[1].Email {
		t.Errorf("Expected %s but got %s", users[1].Email, res.Content[0].Email)
		return
	}
}

func TestDeleteNotExistingRecord(t *testing.T) {
	s := getMockedUserService()
	users := []user.CreateUserRequest{
		{
			Name:     "XXX YYY",
			Email:    "example1@mail.com",
			Password: "!testPass_4251",
		},
		{
			Name:     "XXX ZZZ",
			Email:    "example2@mail.com",
			Password: "!testPass_4251",
		},
	}
	for _, v := range users {
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

func getMockedUserService() *services.UserService {
	mockRepo := mocks.NewMockRepo[entities.User]()
	logger := util.NewLogger(map[string]string{}, os.Stdout)
	hasher := mocks.NewMockHashingManager()
	return services.NewUserService(mockRepo, logger, hasher)
}
