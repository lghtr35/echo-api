package services

import (
	"errors"
	"fmt"
	"reson8-learning-api/managers"
	requests "reson8-learning-api/models/dtos/requests/user"
	responses "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService struct {
	db     *gorm.DB
	hasher managers.HashingManager
	logger *util.Logger
}

func NewUserService(db *gorm.DB, logger *util.Logger, hasher managers.HashingManager) *UserService {
	return &UserService{db: db, logger: logger, hasher: hasher}
}

func (s UserService) GetOne(id string) (entities.User, error) {
	var user entities.User
	s.logger.Debug().Msg(fmt.Sprintf("UserService_GetOne with id: %s", id))
	res := s.db.Preload("Notes").Preload("Languages").First(&user, id)
	if res.Error != nil {
		s.logger.Error().Msg("UserService_GetOne had an error when getting data from db")
		return entities.User{}, res.Error
	}

	return user, nil
}

func (s UserService) FilterAll(request requests.FilterUsersRequest) (responses.PaginationResponse[entities.User], error) {
	var users []entities.User
	s.logger.Debug().Msg(fmt.Sprintf("UserService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	res := q.Find(&users)
	if res.Error != nil {
		s.logger.Error().Msg("UserService_FilterAll had an error when requesting the data from db")
		return responses.PaginationResponse[entities.User]{}, res.Error
	}

	var count int64
	res = q.Count(&count)
	if res.Error != nil {
		s.logger.Error().Msg("UserService_FilterAll had an error when requesting the data from db")
		return responses.PaginationResponse[entities.User]{}, res.Error
	}
	return responses.PaginationResponse[entities.User]{Content: users, Page: request.Page, Size: len(users), TotalCount: int(count)}, nil
}

func (s UserService) buildFilterQuery(request requests.FilterUsersRequest) *gorm.DB {
	q := s.db.Model(&entities.User{})
	s.logger.Debug().Msg("UserService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("UserService filtering IDs")
		q = q.Where(*request.IDs)
	}
	if request.EmailQuery != nil && *request.EmailQuery != "" {
		s.logger.Debug().Msg("UserService filtering Email")
		emailLike := "%" + *request.EmailQuery + "%"
		q = q.Where("email LIKE ?", emailLike)
	}
	if request.NameQuery != nil && *request.NameQuery != "" {
		s.logger.Debug().Msg("UserService filtering Name")
		nameLike := "%" + *request.NameQuery + "%"
		q = q.Where("name LIKE ?", nameLike)
	}
	if request.CourseIDs != nil && len(*request.CourseIDs) > 0 {
		s.logger.Debug().Msg("UserService filtering Course")
		q = q.Where("courseID IN ?", *request.CourseIDs)
	}

	return q.Order("created_at")
}

// TODO do email and password check if not good enough reject creation
func (s UserService) CreateOne(request requests.CreateUserRequest) (entities.User, error) {
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
	res := s.db.Create(&user)
	if res.Error != nil {
		s.logger.Error().Msg("UserService_CreateOne had an error when saving to db")
		return entities.User{}, err
	}
	res = s.db.Save(&user)
	if res.Error != nil {
		s.logger.Error().Msg("UserService_CreateOne had an error when saving to db")
		return entities.User{}, err
	}
	return user, nil
}

func (s UserService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("UserService_DeleteOne has started with given id: %s", id))
	res := s.db.Select(clause.Associations).Delete(&entities.User{}, id)
	if res.Error != nil {
		s.logger.Error().Msg("UserService_DeleteOne had an error when deleting from db")
		return false, res.Error
	}

	return true, nil
}

func (s UserService) UpdateOne(request requests.UpdateUserRequest) (entities.User, error) {
	s.logger.Debug().Msg(fmt.Sprintf("UserService_UpdateOne has started with given id: %s", request.ID))
	var user entities.User
	res := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, request.ID)
	err := res.Error
	if err != nil {
		s.logger.Error().Msg(fmt.Sprintf("UserService_UpdateOne could not find a record with given id: %s", request.ID))
		return entities.User{}, err
	}

	if request.Name != nil && *request.Name != "" {
		s.logger.Debug().Msg(fmt.Sprintf("UserService_UpdateOne updating Name. From: %v => To: %v", user.Name, *request.Name))
		user.Name = *request.Name
	}

	if request.Languages != nil && len(*request.Languages) > 0 {
		s.logger.Debug().Msg("UserService_UpdateOne updating Languages")
		user, err = s.updateLanguagesForUser(user, request)
		if err != nil {
			s.logger.Error().Msg("UserService_UpdateOne had an error when updating Languages for the user")
			return user, err
		}
	}

	if request.Notes != nil && len(*request.Notes) > 0 {
		s.logger.Debug().Msg("UserService_UpdateOne updating Notes")
		user, err = s.updateNotesForUser(user, request)
		if err != nil {
			s.logger.Error().Msg("UserService_UpdateOne had an error when updating Notes for the user")
			return user, err
		}
	}

	res = s.db.Save(&user)
	if res.Error != nil {
		s.logger.Error().Msg("UserService_UpdateOne had an error while trying to save to db")
		return entities.User{}, res.Error
	}
	return user, nil
}

func (s UserService) updateNotesForUser(user entities.User, request requests.UpdateUserRequest) (entities.User, error) {
	if request.Notes == nil || len(*request.Notes) == 0 {
		return entities.User{}, errors.New("argumentErrorNote")
	}

	s.logger.Debug().Msg("UserService_UpdateOne_UpdateNotes started filtering operations")
	onesToAdd, onesToDelete := filterOperations(*request.Notes, s.logger)

	s.logger.Debug().Msg("UserService_UpdateOne_UpdateNotes trying to make operations on db")
	return saveAssociationUpdatesToDb(s.db, user, onesToAdd, onesToDelete, "Notes")
}

func (s UserService) updateLanguagesForUser(user entities.User, request requests.UpdateUserRequest) (entities.User, error) {
	if request.Languages == nil || len(*request.Languages) == 0 {
		return entities.User{}, errors.New("argumentErrorLanguage")
	}

	s.logger.Debug().Msg("UserService_UpdateOne_UpdateLanguages started filtering operations")
	onesToAdd, onesToDelete := filterOperations(*request.Languages, s.logger)

	s.logger.Debug().Msg("UserService_UpdateOne_UpdateLanguages trying to make operations on db")
	return saveAssociationUpdatesToDb(s.db, user, onesToAdd, onesToDelete, "Languages")
}

func (s UserService) MakeAdmin(id string) error {
	var user entities.User
	res := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, id)
	err := res.Error
	if err != nil {
		return err
	}

	res = s.db.Model(&user).Update("role", entities.Admin)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s UserService) MakeNonAdminRole(id string, role uint) error {
	enum := entities.Role(role)
	if enum == entities.Admin {
		return errors.New("argumentErrorRole")
	}
	var user entities.User
	res := s.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, id)
	err := res.Error
	if err != nil {
		return err
	}

	res = s.db.Model(&user).Update("role", enum)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
