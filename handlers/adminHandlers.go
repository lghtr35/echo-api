package handlers

import (
	"echo-api/models/dtos/requests/language"
	"echo-api/models/dtos/requests/note"
	"echo-api/models/dtos/requests/user"
	_ "echo-api/models/dtos/responses/pagination"
	"echo-api/services"
	"echo-api/util"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	logger          *util.Logger
	userService     *services.UserService
	noteService     *services.NoteService
	languageService *services.LanguageService
}

func InitializeAdminHandlers(logger *util.Logger, us *services.UserService, ns *services.NoteService, ls *services.LanguageService) *AdminHandlers {
	return &AdminHandlers{logger: logger, userService: us, noteService: ns, languageService: ls}
}

func (h *AdminHandlers) ConfigureRoutes(api *gin.RouterGroup) {
	api.PATCH("/users/:id/makeadmin", h.MakeUserAdmin)
	//TODO remove these since dont want a backdoor on user data
	api.GET("/users", h.ReadUserWithFilter)
	api.GET("/notes", h.ReadNoteWithFilter)
	api.GET("/languages", h.ReadLanguageWithFilter)
	api.POST("/languages", h.CreateLanguage)
	api.PATCH("/languages", h.UpdateLanguage)
	api.DELETE("/languages/:id", h.DeleteLanguage)
}

// @BasePath /admin

// MakeUserAdmin godoc
// @Summary Promotes a user to admin status.
// @Schemes
// @Description Updates the role of a user to admin using their ID.
// @Security JwtAuth
// @Tags admin
// @Accept json
// @Produce plain
// @Param id path int true "User ID"
// @Success 200 "Successfully promoted to admin"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/users/{id}/makeadmin [patch]
func (h *AdminHandlers) MakeUserAdmin(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.MakeAdmin(id)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "")
}

// ReadUserWithFilterAdmin godoc
// @Summary Reads users based on filter criteria.
// @Schemes
// @Description Retrieves a list of users that match the specified filter criteria.
// @Security JwtAuth
// @Tags admin, users
// @Accept json
// @Produce json
// @Param filter query user.FilterUsersRequest true "Filter parameters"
// @Success 200 {object} pagination.PaginationResponse[entities.User] "Filtered users"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/users [get]
func (h *AdminHandlers) ReadUserWithFilter(c *gin.Context) {
	var request user.FilterUsersRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	users, err := h.userService.FilterAll(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
}

// ReadNoteWithFilterAdmin godoc
// @Summary Reads notes based on filter criteria.
// @Schemes
// @Description Retrieves a list of notes that match the specified filter criteria.
// @Security JwtAuth
// @Tags admin, notes
// @Accept json
// @Produce json
// @Param filter query note.FilterNotesRequest true "Filter parameters"
// @Success 200 {object} pagination.PaginationResponse[entities.Note] "Filtered notes"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/notes [get]
func (h *AdminHandlers) ReadNoteWithFilter(c *gin.Context) {
	var request note.FilterNotesRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	notes, err := h.noteService.FilterAll(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, notes)
}

// ReadLanguageWithFilterAdmin godoc
// @Summary Reads languages based on filter criteria.
// @Schemes
// @Description Retrieves a list of languages that match the specified filter criteria.
// @Security JwtAuth
// @Tags admin, languages
// @Accept json
// @Produce json
// @Param filter query language.FilterLanguagesRequest true "Filter parameters"
// @Success 200 {object} pagination.PaginationResponse[entities.Language] "Filtered languages"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/languages [get]
func (h *AdminHandlers) ReadLanguageWithFilter(c *gin.Context) {
	var request language.FilterLanguagesRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	languages, err := h.languageService.FilterAll(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, languages)
}

// CreateLanguage godoc
// @Summary Creates a new language.
// @Schemes
// @Description Accepts a payload to create a new language and returns the created language ID.
// @Security JwtAuth
// @Tags admin, languages
// @Accept json
// @Produce json
// @Param request body language.CreateLanguageRequest true "Create Language Request"
// @Success 200 {object} map[string]interface{} "Language ID response"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/languages [post]
func (h *AdminHandlers) CreateLanguage(c *gin.Context) {
	var request language.CreateLanguageRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language, err := h.languageService.CreateOne(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"language": language})
}

// DeleteLanguage godoc
// @Summary Deletes a language by ID.
// @Schemes
// @Description Deletes the language associated with the provided ID.
// @Security JwtAuth
// @Tags admin, languages
// @Accept json
// @Produce json
// @Param id path int true "Language ID"
// @Success 200 {object} map[string]interface{} "Deletion success status"
// @Failure 400 {object} string "Bad Request"
// @Failure 404 {object} string "Not Found"
// @Router /admin/languages/{id} [delete]
func (h *AdminHandlers) DeleteLanguage(c *gin.Context) {
	id := c.Param("id")

	ok, err := h.languageService.DeleteOne(id)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !ok {
		h.logger.Err(errors.New("notFoundError"))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"isOk": ok})
}

// UpdateLanguage godoc
// @Summary Updates an existing language.
// @Schemes
// @Description Updates the details of a language based on the provided payload.
// @Security JwtAuth
// @Tags admin, languages
// @Accept json
// @Produce json
// @Param request body language.UpdateLanguageRequest true "Update Language Request"
// @Success 200 {object} entities.Language "Updated language"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /admin/languages [patch]
func (h *AdminHandlers) UpdateLanguage(c *gin.Context) {
	var request language.UpdateLanguageRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	language, err := h.languageService.UpdateOne(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, language)
}
