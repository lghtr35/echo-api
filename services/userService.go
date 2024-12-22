package services

import (
	"errors"
	"fmt"
	"reson8-learning-api/managers"
	requests "reson8-learning-api/models/dtos/requests/user"
	responses "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"

	"gorm.io/gorm/clause"
)

type UserService struct {
	repo   util.Repository[entities.User]
	hasher managers.HashingManager
	logger *util.Logger
}

func NewUserService(repo util.Repository[entities.User], logger *util.Logger, hasher managers.HashingManager) *UserService {
	return &UserService{repo: repo, logger: logger, hasher: hasher}
}

func (s *UserService) GetOne(id string) (entities.User, error) {
	s.logger.Debug().Msg(fmt.Sprintf("UserService_GetOne with id: %s", id))
	res, err := s.repo.First(id, true)
	if err != nil {
		s.logger.Error().Msg("UserService_GetOne had an error when getting data from repo")
		return entities.User{}, err
	}

	return res, nil
}

func (s *UserService) FilterAll(request requests.FilterUsersRequest) (responses.PaginationResponse[entities.User], error) {
	s.logger.Debug().Msg(fmt.Sprintf("UserService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	res, err := q.Find(true)
	if err != nil {
		s.logger.Error().Msg("UserService_FilterAll had an error when requesting the data from repo")
		return responses.PaginationResponse[entities.User]{}, err
	}

	count, err := q.Count()
	if err != nil {
		s.logger.Error().Msg("UserService_FilterAll had an error when requesting the data from repo")
		return responses.PaginationResponse[entities.User]{}, err
	}
	return responses.PaginationResponse[entities.User]{Content: res, Page: request.Page, Size: len(res), TotalCount: int(count)}, nil
}

func (s *UserService) buildFilterQuery(request requests.FilterUsersRequest) util.Repository[entities.User] {
	q := s.repo.Query()
	s.logger.Debug().Msg("*UserService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("*UserService filtering IDs")
		q = q.Where("ID IN ?", *request.IDs)
	}
	if request.EmailQuery != nil && *request.EmailQuery != "" {
		s.logger.Debug().Msg("*UserService filtering Email")
		emailLike := "%" + *request.EmailQuery + "%"
		q = q.Where("email LIKE ?", emailLike)
	}
	if request.NameQuery != nil && *request.NameQuery != "" {
		s.logger.Debug().Msg("*UserService filtering Name")
		nameLike := "%" + *request.NameQuery + "%"
		q = q.Where("name LIKE ?", nameLike)
	}
	if request.CourseIDs != nil && len(*request.CourseIDs) > 0 {
		s.logger.Debug().Msg("*UserService filtering Course")
		q = q.Where("courseID IN ?", *request.CourseIDs)
	}

	return q.Order("created_at")
}

// TODO do email and password check if not good enough reject creation
func (s *UserService) CreateOne(request requests.CreateUserRequest) (entities.User, error) {
	s.logger.Debug().Msg("UserService_CreateOne has started")
	s.logger.Debug().Msg("UserService_CreateOne trying to get password hash")
	hash, err := s.hasher.GetHash(request.Password)
	if err != nil {
		s.logger.Error().Msg("UserService_CreateOne had an error when trying to hash password")
		return entities.User{}, err
	}

	user := entities.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: entities.Password{Value: hash},
		Role:     entities.Customer,
	}
	user, err = s.repo.Create(&user)
	if err != nil {
		s.logger.Error().Msg("UserService_CreateOne had an error when saving to repo")
		return entities.User{}, err
	}
	return user, nil
}

func (s *UserService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("UserService_DeleteOne has started with given id: %s", id))
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error().Msg("UserService_DeleteOne had an error when deleting from repo")
		return false, err
	}

	return true, nil
}

func (s *UserService) UpdateOne(request requests.UpdateUserRequest) (entities.User, error) {
	s.logger.Debug().Msg(fmt.Sprintf("UserService_UpdateOne has started with given id: %s", request.ID))
	user, err := s.repo.Clauses(clause.Locking{Strength: "UPDATE"}).First(request.ID, true)
	if err != nil {
		s.logger.Error().Msg(fmt.Sprintf("UserService_UpdateOne could not find a record with given id: %s", request.ID))
		return entities.User{}, err
	}

	if request.Name != nil && *request.Name != "" {
		s.logger.Debug().Msg(fmt.Sprintf("UserService_UpdateOne updating Name. From: %v => To: %v", user.Name, *request.Name))
		user.Name = *request.Name
	}

	user, err = s.repo.Update(&user)
	if err != nil {
		s.logger.Error().Msg("UserService_UpdateOne had an error while trying to save to repo")
		return entities.User{}, err
	}
	return user, nil
}

func (s *UserService) MakeAdmin(id string) error {
	user, err := s.repo.Clauses(clause.Locking{Strength: "UPDATE"}).First(id, false)
	if err != nil {
		return err
	}
	user.Role = entities.Admin

	_, err = s.repo.Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) MakeNonAdminRole(id string, role uint) error {
	enum := entities.Role(role)
	if enum == entities.Admin {
		return errors.New("argumentErrorRole")
	}
	res, err := s.repo.Clauses(clause.Locking{Strength: "UPDATE"}).First(id, false)
	if err != nil {
		return err
	}
	res.Role = enum

	_, err = s.repo.Update(&res)
	if err != nil {
		return err
	}
	return nil
}
