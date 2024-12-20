package handlers

import "github.com/gin-gonic/gin"

type HandlersBase interface {
	ConfigureRoutes(*gin.RouterGroup)
}
