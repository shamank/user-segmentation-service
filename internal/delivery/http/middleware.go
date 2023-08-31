package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
)

func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

var isSlug validator.Func = func(fl validator.FieldLevel) bool {
	slug := fl.Field().Interface().(string)
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", slug)
	return match
}
