package handlers

import (
	"echo-api/models/dtos/requests/auth"
	"echo-api/models/dtos/requests/user"
	"echo-api/services"
	"echo-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnonymousHandlers struct {
	logger      *util.Logger
	authService *services.AuthService
	userService *services.UserService
}

func InitializeAnonymousHandlers(logger *util.Logger, us *services.UserService, as *services.AuthService) *AnonymousHandlers {
	return &AnonymousHandlers{logger: logger, userService: us, authService: as}
}

func (h *AnonymousHandlers) ConfigureRoutes(api *gin.RouterGroup) {
	api.POST("/login", h.Login)
	api.POST("/register", h.CreateUser)
}

// @BasePath

// Login godoc
// @Summary Authenticates a user and generates a token.
// @Schemes
// @Description Handles user login requests by validating credentials and returning a token.
// @Tags anon, auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{} "Token response"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /login [post]
func (h *AnonymousHandlers) Login(c *gin.Context) {
	var request auth.LoginRequest
	err := c.Bind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, err := h.authService.Login(request)
	if err != nil {
		h.logger.Err(err)
		if err.Error() != "passwordIncorrect" {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"Token": token})
}

// CreateUser godoc
// @Summary Creates a new user.
// @Schemes
// @Description Handles user creation requests by accepting a payload and returning the created user ID.
// @Tags anon, users
// @Accept json
// @Produce json
// @Param request body user.CreateUserRequest true "Create User Request"
// @Success 200 {object} map[string]interface{} "User ID response"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /register [post]
func (h *AnonymousHandlers) CreateUser(c *gin.Context) {
	var request user.CreateUserRequest
	err := c.Bind(&request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := h.userService.CreateOne(request)
	if err != nil {
		h.logger.Err(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"id": id})
}
