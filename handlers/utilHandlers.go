package handlers

import (
	"echo-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type UtilHandlers struct {
	config *util.Configuration
}

func InitializeUtilHandlers(c *util.Configuration) *UtilHandlers {
	return &UtilHandlers{config: c}
}

func (h *UtilHandlers) ConfigureRoutes(api *gin.RouterGroup) {
	api.GET("/", h.GetHome)
	api.GET("/healthcheck", h.GetHealthCheck)
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// @BasePath /api

// GetHome godoc
// @Summary Gets home html
// @Schemes
// @Description serve basic html
// @Tags util
// @Accept json
// @Produce html
// @Success 200
// @Router /api [get]
func (h *UtilHandlers) GetHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":       h.config.Title,
		"swaggerLink": h.config.SwaggerUrl,
		"version":     h.config.Version,
	})
}

// GetHealthCheck godoc
// @Summary Healthcheck
// @Schemes
// @Description Get the status on Config, DB conn, Logging, Services
// @Tags util
// @Accept json
// @Produce json
// @Success 200
// @Fail 500
// @Router /api/healthcheck [get]
func (h *UtilHandlers) GetHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]any{
		"Configuration": true,
		"DB":            true,
		"Logger":        true,
		"Services":      true,
		"Server":        true,
	})
}
