package v1

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, msg string) {

	c.AbortWithStatusJSON(statusCode, errorResponse{Status: "error", Message: msg})
}
