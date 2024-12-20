package internal

import (
	"os"
	"reson8-learning-api/handlers"
	"reson8-learning-api/managers"
	"reson8-learning-api/managers/implementations"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/services"
	"reson8-learning-api/util"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var configuration *util.Configuration
var db *gorm.DB
var logger *util.Logger
var hasher managers.HashingManager
var fileManager managers.FileManager
var promptManager managers.PromptGenManager
var aiCommunicationManager managers.AiCommunicationManager

var authService *services.AuthService
var documentService *services.DocumentService
var languageService *services.LanguageService
var noteService *services.NoteService
var userService *services.UserService
var contextService *services.ContextService
var promptService *services.PromptService

var utilHandlers *handlers.UtilHandlers
var anonymousHandlers *handlers.AnonymousHandlers
var authorizedHandlers *handlers.AuthorizedHandlers
var adminHandlers *handlers.AdminHandlers

func InjectDeps() error {
	var err error
	logger = util.NewLogger(map[string]string{}, os.Stdout)

	configuration, err = util.NewConfiguration(logger)
	if err != nil {
		return err
	}

	hasher, err = implementations.NewBlake3HashingManager(configuration)
	if err != nil {
		return err
	}

	fileManager = implementations.NewOnServerFileManager("~/FileSaveLoc", configuration.SaveLocations)

	promptManager = implementations.NewLocalPromptGenManager(fileManager)

	aiCommunicationManager = implementations.NewOpenAiCommunicationManager(configuration)

	db, err = gorm.Open(postgres.Open(configuration.DbConnectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DoMigrationsIfExists()
	if err != nil {
		return err
	}

	configureServices()

	InitializeHandlers()

	return nil
}

func configureServices() {
	authService = services.NewAuthService(db, hasher, logger, configuration.GetSecretKey())
	documentService = services.NewDocumentService(db, logger, fileManager)
	languageService = services.NewLanguageService(db, logger)
	noteService = services.NewNoteService(db, logger)
	userService = services.NewUserService(db, logger, hasher)
	contextService = services.NewContextService(db, logger)
	promptService = services.NewPromptService(db, logger, promptManager, aiCommunicationManager)
}

func DoMigrationsIfExists() error {
	err := db.AutoMigrate(
		&entities.User{},
		&entities.Document{},
		&entities.Note{},
		&entities.Context{},
		&entities.Prompt{},
		&entities.Password{},
	)
	if err != nil {
		return err
	}
	err = db.Set("gorm:table_options", "CHARSET=utf8mb4").
		AutoMigrate(&entities.Language{})
	if err != nil {
		return err
	}
	return nil
}

func InitializeHandlers() {
	utilHandlers = handlers.InitializeUtilHandlers(configuration)
	anonymousHandlers = handlers.InitializeAnonymousHandlers(logger, userService, authService)
	authorizedHandlers = handlers.InitializeAuthorizedHandlers(logger, userService, authService, noteService, languageService, documentService, contextService, promptService)
	adminHandlers = handlers.InitializeAdminHandlers(logger, userService, noteService, languageService)
}

func MapEnpoints(g *gin.Engine) {
	api := g.Group("/api")
	api.Use(authService.CORSMiddleware())
	{
		utilHandlers.ConfigureRoutes(api)
		v1 := api.Group("/v1")
		{
			anonymousHandlers.ConfigureRoutes(v1)
			authorized := v1.Group("/")
			authorized.Use(authService.AuthMiddleware())
			{
				authorizedHandlers.ConfigureRoutes(authorized)
				admin := authorized.Group("/admin")
				admin.Use(authService.AdminMiddleware())
				{
					adminHandlers.ConfigureRoutes(admin)
				}
			}
		}
	}
}

func GetConfiguration() *util.Configuration {
	return configuration
}

func GetLogger() *util.Logger {
	return logger
}

func GetHashingManager() *managers.HashingManager {
	return &hasher
}

func GetAuthService() *services.AuthService {
	return authService
}

func GetDocumentService() *services.DocumentService {
	return documentService
}

func GetLanguageService() *services.LanguageService {
	return languageService
}
func GetNoteService() *services.NoteService {
	return noteService
}
func GetUserService() *services.UserService {
	return userService
}

func GetContextService() *services.ContextService {
	return contextService
}

func GetPromptService() *services.PromptService {
	return promptService
}
