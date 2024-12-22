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

var noteRepository *util.GormRepository[entities.Note]
var userRepository *util.GormRepository[entities.User]
var documentRepository *util.GormRepository[entities.Document]
var languageRepository *util.GormRepository[entities.Language]
var contextRepository *util.GormRepository[entities.Context]
var promptRepository *util.GormRepository[entities.Prompt]

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

	initializeRepositories()

	configureServices()

	initializeHandlers()

	return nil
}

func initializeRepositories() {
	noteRepository = util.NewGormRepository[entities.Note](db, []string{"Documents"})
	documentRepository = util.NewGormRepository[entities.Document](db, []string{})
	languageRepository = util.NewGormRepository[entities.Language](db, []string{"Notes", "Contexts"})
	userRepository = util.NewGormRepository[entities.User](db, []string{"Contexts", "Documents", "Notes", "Languages"})
	contextRepository = util.NewGormRepository[entities.Context](db, []string{"Notes", "Prompts", "Documents"})
	promptRepository = util.NewGormRepository[entities.Prompt](db, []string{})
}

func configureServices() {
	authService = services.NewAuthService(db, hasher, logger, configuration.GetSecretKey())
	documentService = services.NewDocumentService(documentRepository, logger, fileManager)
	languageService = services.NewLanguageService(languageRepository, logger)
	noteService = services.NewNoteService(noteRepository, logger)
	userService = services.NewUserService(userRepository, logger, hasher)
	contextService = services.NewContextService(contextRepository, logger)
	promptService = services.NewPromptService(promptRepository, logger, promptManager, aiCommunicationManager)
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

func initializeHandlers() {
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
