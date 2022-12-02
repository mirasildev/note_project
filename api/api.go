package api

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/mirasildev/note_project/api/v1"
	"github.com/mirasildev/note_project/config"
	"github.com/mirasildev/note_project/storage"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/mirasildev/note_project/api/docs" // for swagger
)

type RouterOptions struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

// / @title           Swagger for note api
// @version         1.0
// @description     This is a note taker service api.
// @host      localhost:8000
// @BasePath  /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:      opt.Cfg,
		Storage:  opt.Storage,
		InMemory: opt.InMemory,
	})

	router.Static("/media", "./media")

	apiV1 := router.Group("/v1")

	apiV1.POST("/users", handlerV1.AuthMiddleware, handlerV1.CreateUser)
	apiV1.GET("/users/:id", handlerV1.AuthMiddleware, handlerV1.GetUser)
	apiV1.GET("/users", handlerV1.GetAllUsers)
	apiV1.PUT("/users/:id", handlerV1.UpdateUser)
	apiV1.DELETE("/users/:id", handlerV1.DeleteUser)

	apiV1.POST("/notes", handlerV1.AuthMiddleware, handlerV1.CreateNote)
	apiV1.GET("/notes/:id", handlerV1.GetNote)
	apiV1.GET("/notes", handlerV1.GetAllNotes)
	apiV1.PUT("/notes/:id", handlerV1.UpdateNote)
	apiV1.DELETE("/notes/:id", handlerV1.DeleteNote)

	apiV1.POST("/auth/register", handlerV1.Register)
	apiV1.POST("/auth/login", handlerV1.Login)
	apiV1.POST("/auth/verify", handlerV1.Verify)

	apiV1.POST("/file-upload", handlerV1.AuthMiddleware, handlerV1.UploadFile)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
