package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/shamank/user-segmentation-service/docs"
	v1 "github.com/shamank/user-segmentation-service/internal/delivery/http/v1"
	"github.com/shamank/user-segmentation-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
)

type Handler struct {
	services *service.Services
	logger   *slog.Logger

	basePath string
}

func NewHandler(services *service.Services, logger *slog.Logger, basePath string) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
		basePath: basePath,
	}
}

func (h *Handler) InitRoute() *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.NewHandlerV1(h.services, h.logger, h.basePath)
	router.Use(CORS)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("slug", isSlug)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/file/:filename", func(c *gin.Context) {
		filename := c.Param("filename")

		c.File("./assets/user_logs/" + filename)
	})

	// TODO: init handler v1

	api := router.Group("/api")
	{
		handlerV1.InitRouteV1(api)
	}

	return router
}
