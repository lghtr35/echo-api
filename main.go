package main

import (
	"strconv"

	"echo-api/internal"

	_ "echo-api/docs"

	"github.com/gin-gonic/gin"
)

func Configure() *gin.Engine {
	err := internal.InjectDeps()
	if err != nil {
		logger := internal.GetLogger()
		logger.Fatal().Err(err).Msg("Error occurred while injecting dependencies")
	}
	g := gin.New()
	g.LoadHTMLGlob("static/html/*")
	internal.MapEnpoints(g)

	return g
}

func Start(g *gin.Engine, port int) {
	address := ":" + strconv.Itoa(port)
	g.Run(address)
}

// @title           LanguHelp API
// @version         0.0.1
// @description     AI assisted learning, studying and working tool.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url	   https://www.serdilcakmak.com
// @contact.email  serdilcakmak@gmail.com

// @license.name  Copyright of Serdil Cagin Cakmak
// @license.url

// @host      localhost:11242
// @BasePath  /api/v1

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization
// @description Bearer

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	g := Configure()
	Start(g, 11242)
}
