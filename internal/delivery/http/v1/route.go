package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/user-segmentation-service/internal/service"
	"log/slog"
)

type HandlerV1 struct {
	services *service.Services
	logger   *slog.Logger
	basePath string
}

func NewHandlerV1(services *service.Services, logger *slog.Logger, basePath string) *HandlerV1 {

	return &HandlerV1{
		services: services,
		logger:   logger,
		basePath: basePath,
	}
}

func (h *HandlerV1) InitRouteV1(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initSegmentRoute(v1)
		h.initUserRoute(v1)
	}
}
