package services

import (
	"errors"
	"fmt"
	managers "reson8-learning-api/managers"
	documentRequest "reson8-learning-api/models/dtos/requests/document"
	documentResponse "reson8-learning-api/models/dtos/responses/document"
	responses "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"
	"strings"
	"sync"
)

type DocumentService struct {
	repo        util.Repository[entities.Document]
	logger      *util.Logger
	fileManager managers.FileManager
}

func NewDocumentService(repo util.Repository[entities.Document], logger *util.Logger, manager managers.FileManager) *DocumentService {
	return &DocumentService{repo, logger, manager}
}

func (s *DocumentService) CheckIfBelongsToUser(id string, userID string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("DocumentService_CheckIfBelongsToUser with id: %s for user: %s", id, userID))
	document, err := s.GetOne(id)
	if err != nil {
		s.logger.Error().Msg("DocumentService_CheckIfBelongsToUser had an error when getting from repo")
		return false, err
	}

	return userID == document.UserID, nil
}

func (s *DocumentService) GetOne(id string) (documentResponse.DocumentWrapped, error) {
	s.logger.Debug().Msg(fmt.Sprintf("DocumentService_GetOne with id: %s", id))
	document, err := s.repo.First(id, false)
	if err != nil {
		s.logger.Error().Msg("DocumentService_GetOne had an error when getting from repo")
		return documentResponse.DocumentWrapped{}, err
	}

	return s.mapOneToDocumentWrapped(document), nil
}

func (s *DocumentService) FilterAll(request documentRequest.FilterDocumentsRequest) (responses.PaginationResponse[documentResponse.DocumentWrapped], error) {
	s.logger.Debug().Msg(fmt.Sprintf("DocumentService_FilterAll on page: %d with size: %d", request.Page, request.Size))
	offset := request.CalculateOffset()

	q := s.buildFilterQuery(request)
	q.Offset(int(offset)).Limit(int(request.Size))
	docs, err := q.Find(false)
	if err != nil {
		s.logger.Error().Msg("DocumentService_FilterAll had an error when requesting the data from repo")
		return responses.PaginationResponse[documentResponse.DocumentWrapped]{}, err
	}

	count, err := q.Count()
	if err != nil {
		s.logger.Error().Msg("DocumentService_FilterAll had an error when requesting the data from repo")
		return responses.PaginationResponse[documentResponse.DocumentWrapped]{}, err
	}
	return responses.PaginationResponse[documentResponse.DocumentWrapped]{Page: request.Page, Size: len(docs), TotalCount: int(count), Content: s.mapToDocumentWrapped(docs)}, nil
}

func (s *DocumentService) buildFilterQuery(request documentRequest.FilterDocumentsRequest) util.Repository[entities.Document] {
	q := s.repo.Query()
	s.logger.Debug().Msg("DocumentService started to build Filter query")

	if request.IDs != nil && len(*request.IDs) > 0 {
		s.logger.Debug().Msg("DocumentService filtering IDs")
		q = q.Where("ID IN ?", *request.IDs)
	}

	if request.Name != nil && len(*request.Name) > 0 {
		s.logger.Debug().Msg("DocumentService filtering Name")
		nameLike := "%" + *request.Name + "%"
		q = q.Where("name LIKE ?", nameLike)
	}

	if request.NoteIDs != nil && len(*request.NoteIDs) > 0 {
		s.logger.Debug().Msg("DocumentService filtering documents")
		q = q.Where("noteID in ?", *request.NoteIDs)
	}

	if request.Location != nil && len(*request.Location) > 0 {
		s.logger.Debug().Msg("DocumentService filtering documents")
		q = q.Where("location = ?", *request.Location)
	}

	if request.Extension != nil && len(*request.Extension) > 0 {
		s.logger.Debug().Msg("DocumentService filtering Courses")
		q = q.Where("extension = ?", *request.Extension)
	}

	if request.ContextIDs != nil && len(*request.ContextIDs) > 0 {
		s.logger.Debug().Msg("DocumentService filtering ContextID")
		q = q.Where("contextID IN ?", *request.ContextIDs)
	}

	return q.Order("created_at")
}

func (s *DocumentService) CreateBulkFromMultipart(request documentRequest.CreateDocumentsMultipartRequest) ([]entities.Document, error) {
	if len(request.Files) == 0 {
		return nil, errors.New("argumentErrorMissing")
	}
	s.logger.Debug().Msg("DocumentService_CreateBulkFromMultipart has started")
	count := len(request.Files)
	res := make([]entities.Document, count)
	var wg sync.WaitGroup
	resch := make(chan entities.Document)
	errch := make(chan error)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go s.concurrentCreateOneFromMultipart(&wg, resch, errch, documentRequest.CreateDocumentMultipartRequest{File: request.Files[i], CreateDocumentRequestBase: request.CreateDocumentRequestBase})
	}
	wg.Wait()
	close(resch)
	close(errch)
	for err := range errch {
		if err != nil {
			return nil, err
		}
	}
	i := 0
	for result := range resch {
		res[i] = result
		i++
	}
	return res, nil
}

func (s *DocumentService) CreateOneFromMultipart(request documentRequest.CreateDocumentMultipartRequest) (entities.Document, error) {
	if request.File == nil {
		return entities.Document{}, errors.New("argumentErrorMissing")
	}
	s.logger.Debug().Msg("DocumentService_CreateOneFromMultipart has started")
	err := s.saveMultipartFile(request)
	if err != nil {
		return entities.Document{}, err
	}
	name := request.File.Filename
	extension := s.getFileExtension(name)
	document := entities.Document{
		Name:            name,
		Location:        request.Location,
		Extension:       extension,
		UserID:          request.UserID,
		IsReadableByAll: request.IsReadableByAll,
	}
	if request.EntityType != nil && request.EntityID != nil && *request.EntityType != "" && *request.EntityID != "" {
		document, err = s.addDocumentEntityRelation(document, *request.EntityType, *request.EntityID)
		if err != nil {
			return entities.Document{}, err
		}
	}
	document, err = s.repo.Create(&document)
	if err != nil {
		s.logger.Error().Msg("DocumentService_CreateOneFromMultipart had an error when saving to repo")
		return entities.Document{}, err
	}
	return document, nil
}

func (s *DocumentService) DeleteOne(id string) (bool, error) {
	s.logger.Debug().Msg(fmt.Sprintf("DocumentService_DeleteOne has started with given id: %s", id))
	document, err := s.repo.First(id, false)
	if err != nil {
		s.logger.Error().Msg("DocumentService_DeleteOne had an error when getting from repo")
		return false, err
	}

	err = s.fileManager.DeleteFile(document.Location, document.Name)
	if err != nil {
		return false, err
	}
	err = s.repo.Delete(id)
	if err != nil {
		s.logger.Error().Msg("DocumentService_DeleteOne had an error when deleting from repo")
		return false, err
	}

	return true, nil
}

func (s *DocumentService) saveMultipartFile(request documentRequest.CreateDocumentMultipartRequest) error {
	filename := request.File.Filename
	size := request.File.Size
	f, err := request.File.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	buffSize := int64(32768)
	offset := int64(0)
	buffer := make([]byte, buffSize)
	for offset < size {
		countRead, err := f.ReadAt(buffer, offset)
		if err != nil {
			return err
		}
		countWritten, err := s.fileManager.SaveFile(request.Location, filename, buffer, managers.FileOpeningOptions{StartPoint: managers.CUSTOM, Offset: uint64(offset)})
		if err != nil {
			return err
		} else if countRead != countWritten {
			return errors.New("ioErrorReadWriteMismatch")
		}
		offset += int64(countWritten)
	}

	return nil
}

func (s *DocumentService) getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return parts[len(parts)-1]
}

func (s *DocumentService) concurrentCreateOneFromMultipart(wg *sync.WaitGroup, resch chan entities.Document, errch chan error, request documentRequest.CreateDocumentMultipartRequest) {
	defer wg.Done()
	res, err := s.CreateOneFromMultipart(request)
	resch <- res
	errch <- err
}

func (s *DocumentService) addDocumentEntityRelation(document entities.Document, entityType string, id string) (entities.Document, error) {
	if strings.ToLower(entityType) == "note" {
		document.NoteID = &id
	} else {
		return document, errors.New("notImplementedOwnerType")
	}

	return document, nil
}

func (s *DocumentService) mapOneToDocumentWrapped(doc entities.Document) documentResponse.DocumentWrapped {
	return documentResponse.DocumentWrapped{Document: doc, Path: s.fileManager.GetFullPath(doc.Location, doc.Name)}
}

func (s *DocumentService) mapToDocumentWrapped(docs []entities.Document) []documentResponse.DocumentWrapped {
	res := make([]documentResponse.DocumentWrapped, len(docs))
	for i, v := range docs {
		res[i] = s.mapOneToDocumentWrapped(v)
	}

	return res
}
