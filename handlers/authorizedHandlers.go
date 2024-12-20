package handlers

import (
	"errors"
	"net/http"
	"reson8-learning-api/models/dtos/requests/context"
	"reson8-learning-api/models/dtos/requests/document"
	"reson8-learning-api/models/dtos/requests/language"
	"reson8-learning-api/models/dtos/requests/note"
	"reson8-learning-api/models/dtos/requests/prompt"
	"reson8-learning-api/models/dtos/requests/user"
	_ "reson8-learning-api/models/dtos/responses/pagination"
	"reson8-learning-api/services"
	"reson8-learning-api/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthorizedHandlers struct {
	logger          *util.Logger
	authService     *services.AuthService
	userService     *services.UserService
	noteService     *services.NoteService
	languageService *services.LanguageService
	documentService *services.DocumentService
	contextService  *services.ContextService
	promptService   *services.PromptService
}

func InitializeAuthorizedHandlers(logger *util.Logger, us *services.UserService, as *services.AuthService, ns *services.NoteService, ls *services.LanguageService, ds *services.DocumentService, cs *services.ContextService, ps *services.PromptService) *AuthorizedHandlers {
	return &AuthorizedHandlers{logger: logger, userService: us, authService: as, noteService: ns, languageService: ls, documentService: ds, contextService: cs, promptService: ps}
}

func (h *AuthorizedHandlers) ConfigureRoutes(api *gin.RouterGroup) {
	api.GET("/users/:id", h.ReadUserWithID)
	api.PATCH("/users", h.UpdateUser)
	api.DELETE("users/:id", h.DeleteUser)
	api.PATCH("/users/:id/:role", h.MakeUserNonAdmin)

	api.POST("/notes", h.CreateNote)
	api.GET("/notes/:id", h.ReadNoteWithID)
	api.GET("/notes", h.ReadNoteWithFilter)
	api.PATCH("/notes", h.UpdateNote)
	api.PATCH("/notes/document", h.CreateNoteDocuments)
	api.DELETE("/notes/:id", h.DeleteNote)

	api.GET("/languages/:id", h.ReadLanguageWithID)
	api.GET("/languages", h.ReadLanguageWithFilter)

	api.POST("/documents", h.CreateUserDocument)
	api.POST("/documents/bulk", h.CreateUserDocumentBulk)
	api.GET("/documents/:id", h.ReadUserDocumentWithID)
	api.GET("/documents", h.ReadUserDocumentWithFilter)
	api.DELETE("/documents/:id", h.DeleteDocument)

	api.POST("/contexts", h.CreateContext)
	api.POST("/contexts/:id", h.DeleteContext)
}

// @BasePath /admin

// ReadUserWithID godoc
// @Summary Retrieves a user by ID.
// @Schemes
// @Description Fetches details of a user based on the provided ID. Only the user themselves or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.User "User details"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{id} [get]
func (h *AuthorizedHandlers) ReadUserWithID(c *gin.Context) {
	id := c.Param("id")
	if !h.isUserActingOnSelf(c, id, "User") {
		return
	}

	user, err := h.userService.GetOne(id)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Deletes a user by ID.
// @Schemes
// @Description Deletes the user associated with the provided ID. Only the user themselves or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Deletion success status"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{id} [delete]
func (h *AuthorizedHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if !h.isUserActingOnSelf(c, id, "User") {
		return
	}

	ok, err := h.userService.DeleteOne(id)
	if err != nil || !ok {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"isOk": ok})
}

// UpdateUser godoc
// @Summary Updates user information.
// @Schemes
// @Description Updates the details of a user based on the provided payload. Only the user themselves or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, users
// @Accept json
// @Produce json
// @Param request body user.UpdateUserRequest true "Update User Request"
// @Success 200 {object} entities.User "Updated user details"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users [patch]
func (h *AuthorizedHandlers) UpdateUser(c *gin.Context) {
	var request user.UpdateUserRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if !h.isUserActingOnSelf(c, request.ID, "User") {
		return
	}

	user, err := h.userService.UpdateOne(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

// MakeUserNonAdmin godoc
// @Summary Changes a user’s role to a non-admin role.
// @Schemes
// @Description Updates a user’s role to a non-admin role using their ID and the new role value. Only the user themselves or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, users
// @Accept json
// @Produce plain
// @Param id path int true "User ID"
// @Param role path int true "New role ID"
// @Success 200 "Successfully changed role"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{id}/make-non-admin/{role} [patch]
func (h *AuthorizedHandlers) MakeUserNonAdmin(c *gin.Context) {
	id := c.Param("id")
	if !h.isUserActingOnSelf(c, id, "User") {
		return
	}
	role, err := strconv.ParseUint(c.Param("role"), 10, 32)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = h.userService.MakeNonAdminRole(id, uint(role))
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "")
}

// ReadNoteWithID godoc
// @Summary Retrieves a note by ID.
// @Schemes
// @Description Fetches the details of a specific note based on its ID. Only the owner of the note or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, notes
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} entities.Note "Note details"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /notes/{id} [get]
func (h *AuthorizedHandlers) ReadNoteWithID(c *gin.Context) {
	id := c.Param("id")
	if !h.isUserActingOnSelf(c, id, "Note") {
		return
	}

	note, err := h.noteService.GetOne(id)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, note)
}

// ReadNoteWithFilter godoc
// @Summary Retrieves notes based on filter criteria.
// @Schemes
// @Description Fetches a list of notes that match the specified filter criteria. Notes will only be retrieved for the authorized user.
// @Security JwtAuth
// @Tags authorized, notes
// @Accept json
// @Produce json
// @Param filter query note.FilterNotesRequest true "Filter parameters"
// @Success 200 {object} pagination.PaginationResponse[entities.Note] "Filtered notes"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /notes [get]
func (h *AuthorizedHandlers) ReadNoteWithFilter(c *gin.Context) {
	var request note.FilterNotesRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	request.UserIDs = &[]string{id}

	users, err := h.noteService.FilterAll(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateNote godoc
// @Summary Creates a new note.
// @Schemes
// @Description Accepts a payload to create a new note and associates it with the authenticated user if a user ID is not provided.
// @Security JwtAuth
// @Tags authorized, notes
// @Accept json
// @Produce json
// @Param request body note.CreateNoteRequest true "Create Note Request"
// @Success 200 {object} map[string]interface{} "Note ID response"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /notes [post]
func (h *AuthorizedHandlers) CreateNote(c *gin.Context) {
	var request note.CreateNoteRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	request.UserID = &userID

	note, err := h.noteService.CreateOne(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = h.sendPrompt(note.ContextID, note.ID, note)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{"note": note, "aiError": err.Error()})
	}
	c.JSON(http.StatusOK, map[string]any{"note": note})
}

// DeleteNote godoc
// @Summary Deletes a note by ID.
// @Schemes
// @Description Deletes the note associated with the provided ID. Only the owner of the note or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, notes
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} map[string]interface{} "Deletion success status"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /notes/{id} [delete]
func (h *AuthorizedHandlers) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	if !h.isUserActingOnSelf(c, id, "Note") {
		return
	}

	ok, err := h.noteService.DeleteOne(id)
	if err != nil || !ok {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = h.deletePrompt("", id)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{"isOk": ok, "aiError": err.Error()})
	}
	c.JSON(http.StatusOK, map[string]any{"isOk": ok})
}

// UpdateNote godoc
// @Summary Updates a note.
// @Schemes
// @Description Updates the details of a specific note based on the provided payload. Only the owner of the note or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, notes
// @Accept json
// @Produce json
// @Param request body note.UpdateNoteRequest true "Update Note Request"
// @Success 200 {object} entities.Note "Updated note details"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /notes [patch]
func (h *AuthorizedHandlers) UpdateNote(c *gin.Context) {
	var request note.UpdateNoteRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if !h.isUserActingOnSelf(c, request.ID, "Note") {
		return
	}
	request.UserID = nil

	note, err := h.noteService.UpdateOne(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = h.updatePrompt(note.ContextID, note.ID, note)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{"value": note, "aiError": err.Error()})
	}
	c.JSON(http.StatusOK, note)
}

// ReadLanguageWithID godoc
// @Summary Retrieves a language by ID.
// @Schemes
// @Description Fetches the details of a specific language based on its ID.
// @Security JwtAuth
// @Tags authorized, languages
// @Accept json
// @Produce json
// @Param id path int true "Language ID"
// @Success 200 {object} entities.Language "Language details"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /languages/{id} [get]
func (h *AuthorizedHandlers) ReadLanguageWithID(c *gin.Context) {
	id := c.Param("id")

	language, err := h.languageService.GetOne(id)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, language)
}

// ReadLanguageWithFilter godoc
// @Summary Retrieves languages based on filter criteria.
// @Schemes
// @Description Fetches a list of languages that match the specified filter criteria.
// @Security JwtAuth
// @Tags authorized, languages
// @Accept json
// @Produce json
// @Param filter query language.FilterLanguagesRequest true "Filter parameters"
// @Success 200 {object} pagination.PaginationResponse[entities.Language] "Filtered languages"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /languages [get]
func (h *AuthorizedHandlers) ReadLanguageWithFilter(c *gin.Context) {
	var request language.FilterLanguagesRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	users, err := h.languageService.FilterAll(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
}

// ReadUserDocumentWithID godoc
// @Summary Retrieves a user document by ID.
// @Schemes
// @Description Fetches the details of a specific user document based on its ID. Only the owner of the document or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, documents
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} document.DocumentWrapped "Document details"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /documents/{id} [get]
func (h *AuthorizedHandlers) ReadUserDocumentWithID(c *gin.Context) {
	id := c.Param("id")

	if !h.isUserActingOnSelf(c, id, "Document") {
		return
	}

	document, err := h.documentService.GetOne(id)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, document)
}

// ReadUserDocumentWithFilter godoc
// @Summary Retrieves user documents based on filter criteria.
// @Schemes
// @Description Fetches a list of user documents that match the specified filter criteria.
// @Security JwtAuth
// @Tags authorized, documents
// @Accept json
// @Produce json
// @Param filter query document.FilterDocumentsRequest true "Filter parameters"
// @Success 200 {object} pagination.PaginationResponse[entities.Document] "Filtered documents"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /documents [get]
func (h *AuthorizedHandlers) ReadUserDocumentWithFilter(c *gin.Context) {
	var request document.FilterDocumentsRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	request.UserIDs = &[]string{id}

	docs, err := h.documentService.FilterAll(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, docs)
}

// CreateUserDocument godoc
// @Summary Creates a user document from multipart data.
// @Schemes
// @Description Creates a new document associated with the authenticated user, using multipart data for file upload.
// @Security JwtAuth
// @Tags authorized, documents
// @Accept multipart/form-data
// @Produce json
// @Param request body document.CreateDocumentMultipartRequest true "Create Document Request"
// @Success 200 {object} map[string]interface{} "Document ID"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /documents [post]
func (h *AuthorizedHandlers) CreateUserDocument(c *gin.Context) {
	var request document.CreateDocumentMultipartRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	request.UserID = userID
	request.IsReadableByAll = false

	doc, err := h.documentService.CreateOneFromMultipart(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = h.sendPrompt(doc.ContextID, doc.ID, doc)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{"doc": doc, "aiError": err.Error()})
	}
	c.JSON(http.StatusOK, map[string]any{"doc": doc})
}

// CreateUserDocumentBulk godoc
// @Summary Creates multiple user documents from multipart data.
// @Schemes
// @Description Creates multiple documents for the authenticated user, using multipart data for file uploads.
// @Security JwtAuth
// @Tags authorized, documents
// @Accept multipart/form-data
// @Produce json
// @Param request body document.CreateDocumentsMultipartRequest true "Create Bulk Documents Request"
// @Success 200 {object} map[string]interface{} "Bulk Document IDs"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /documents/bulk [post]
func (h *AuthorizedHandlers) CreateUserDocumentBulk(c *gin.Context) {
	var request document.CreateDocumentsMultipartRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	request.UserID = userID
	request.IsReadableByAll = false

	docs, err := h.documentService.CreateBulkFromMultipart(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ids := make([]string, len(docs))
	errs := make(map[string]string)
	for i, doc := range docs {
		_, err = h.sendPrompt(doc.ContextID, doc.ID, doc)
		if err != nil {
			errs[doc.ID] = err.Error()
		}
		ids[i] = doc.ID
	}

	if len(errs) == 0 {
		c.JSON(http.StatusOK, map[string]any{"ids": ids})
	} else {
		c.JSON(http.StatusOK, map[string]any{"ids": ids, "aiErrors": errs})
	}
}

// DeleteUserDocument godoc
// @Summary Deletes a user document by ID.
// @Schemes
// @Description Deletes a specific document for the authenticated user based on document ID.
// @Security JwtAuth
// @Tags authorized, documents
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} map[string]interface{} "Deletion status"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /documents/{id} [delete]
func (h *AuthorizedHandlers) DeleteDocument(c *gin.Context) {
	id := c.Param("id")

	if !h.isUserActingOnSelf(c, id, "Document") {
		return
	}

	ok, err := h.documentService.DeleteOne(id)
	if err != nil || !ok {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = h.deletePrompt("", id)
	if err != nil {
		c.JSON(http.StatusOK, map[string]any{"isOk": ok, "aiError": err.Error()})
	}
	c.JSON(http.StatusOK, map[string]any{"isOk": ok})
}

// CreateNoteDocuments godoc
// @Summary Creates note-related documents from multipart data.
// @Schemes
// @Description Creates documents linked to a specific note, using multipart data for file uploads. The documents will be associated with the note ID.
// @Security JwtAuth
// @Tags authorized, documents, notes
// @Accept multipart/form-data
// @Produce json
// @Param request body document.CreateNoteDocumentsRequest true "Create Note Documents Request"
// @Success 200 {object} map[string]interface{} "Created document IDs"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /documents/notes [post]
func (h *AuthorizedHandlers) CreateNoteDocuments(c *gin.Context) {
	var request document.CreateNoteDocumentsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	request.UserID = userID
	request.IsReadableByAll = false
	*request.EntityType = "Note"
	request.EntityID = &request.NoteID

	docs, err := h.documentService.CreateBulkFromMultipart(request.CreateDocumentsMultipartRequest)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ids := make([]string, len(docs))
	errs := make(map[string]string)
	for i, doc := range docs {
		_, err = h.sendPrompt(doc.ContextID, doc.ID, doc)
		if err != nil {
			errs[doc.ID] = err.Error()
		}
		ids[i] = doc.ID
	}

	if len(errs) == 0 {
		c.JSON(http.StatusOK, map[string]any{"ids": ids})
	} else {
		c.JSON(http.StatusOK, map[string]any{"ids": ids, "aiErrors": errs})
	}
}

// CreateContext godoc
// @Summary Creates a new context.
// @Schemes
// @Description Accepts a payload to create a new context and associates it with the authenticated user if a user ID is not provided.
// @Security JwtAuth
// @Tags authorized, contexts
// @Accept json
// @Produce json
// @Param request body context.CreateContextRequest true "Create Context Request"
// @Success 200 {object} map[string]interface{} "Context ID response"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /contexts [post]
func (h *AuthorizedHandlers) CreateContext(c *gin.Context) {
	var request context.CreateContextRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	request.UserID = userID

	context, err := h.contextService.CreateOne(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"context": context})
}

// DeleteContext godoc
// @Summary Deletes a context by ID.
// @Schemes
// @Description Deletes the context associated with the provided ID. Only the owner of the context or authorized actions are permitted.
// @Security JwtAuth
// @Tags authorized, contexts
// @Accept json
// @Produce json
// @Param id path int true "Context ID"
// @Success 200 {object} map[string]interface{} "Deletion success status"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /contexts/{id} [delete]
func (h *AuthorizedHandlers) DeleteContext(c *gin.Context) {
	id := c.Param("id")

	if !h.isUserActingOnSelf(c, id, "Context") {
		return
	}

	ok, err := h.contextService.DeleteOne(id)
	if err != nil || !ok {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"isOk": ok})
}

func (h *AuthorizedHandlers) getUserIDFromJwt(c *gin.Context) (string, error) {
	parts := strings.Split(c.Request.Header.Get("Authorization"), " ")
	id, err := h.authService.GetUserIDFromToken(parts[len(parts)-1])
	if err != nil {
		return "", err
	}

	return id, nil
}

func (h *AuthorizedHandlers) isUserActingOnSelf(c *gin.Context, entityID string, entityName string) bool {
	userID, err := h.getUserIDFromJwt(c)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return false
	}
	var ok bool
	switch strings.ToLower(entityName) {
	case "context":
		ok, err = h.contextService.CheckIfBelongsToUser(entityID, userID)
	case "document":
		ok, err = h.documentService.CheckIfBelongsToUser(entityID, userID)
	case "user":
		err = nil
		ok = entityID == userID
	case "note":
		ok, err = h.noteService.CheckIfBelongsToUser(entityID, userID)
	default:
		return true
	}
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return false
	}
	if !ok {
		h.logger.Err(errors.New("authorizationErrorUnauthorizedForContent"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return false
	}

	return true
}

func (h *AuthorizedHandlers) sendPrompt(contextID string, entityID string, val any) (string, error) {

	req := prompt.CreatePromptRequest{
		ContextID: contextID,
		Value:     val,
		EntityID:  entityID,
	}
	p, err := h.promptService.GenerateAndSendPrompt(req)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

func (h *AuthorizedHandlers) updatePrompt(contextID string, entityID string, val any) (string, error) {
	req := prompt.UpdatePromptRequest{
		Value:     val,
		EntityID:  entityID,
		ContextID: contextID,
	}
	p, err := h.promptService.UpdatePrompt(req)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

func (h *AuthorizedHandlers) deletePrompt(contextID string, entityID string) error {
	req := prompt.FindPromptByEntityAndContextRequest{
		EntityID:  entityID,
		ContextID: contextID,
	}
	found, err := h.promptService.FindByEntityAndContext(req)
	if err != nil {
		return err
	}
	_, err = h.promptService.DeleteAndSend(found.ID)
	return err
}
