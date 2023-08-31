package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/user-segmentation-service/pkg/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (h *HandlerV1) initUserRoute(api *gin.RouterGroup) {
	user := api.Group("/user")
	{
		user.GET("", h.getUser)
		user.POST("", h.createUser)

		user.GET("/segments", h.getUserSegments)
		user.POST("/segments", h.addUserToSegments)
		user.DELETE("/segments", h.removeUserFromSegments)

		user.POST("/history", h.getUserHistory)
	}
}

// @Summary Get Profile
// @Tags users
// @Description get user profile
// @ModuleID userGetProfile
// @Accept  json
// @Produce  json
// @Param username path string true "username"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/{username} [get]
func (h *HandlerV1) getUser(c *gin.Context) {

}

type createUserInput struct {
	Username string `json:"username" binding:"required"`
}

// @Summary Create User
// @Tags users
// @Description create new user
// @ModuleID userCreateUser
// @Accept  json
// @Produce  json
// @Param input body createUserInput true "create user"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users [post]
func (h *HandlerV1) createUser(c *gin.Context) {
	const op = "http.v1.user.createUser"
	logger := h.logger.With(slog.String("op", op))

	var input createUserInput

	if err := c.BindJSON(&input); err != nil {
		logger.Error("problem with bind json", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	Id, err := h.services.User.CreateUser(input.Username)
	if err != nil {
		logger.Error("problem createElement", sl.Err(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return

	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"id":     strconv.Itoa(Id),
	})
}

type getUserSegmentsResponse struct {
	Status string   `json:"status"`
	Slugs  []string `json:"slugs"`
}

// @Summary Get User Segments
// @Tags users
// @Description get user segments
// @ModuleID getUserSegments
// @Accept  json
// @Produce  json
// @Param user_id query string true "user id"
// @Success 200 {object} getUserSegmentsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/segments [get]
func (h *HandlerV1) getUserSegments(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "uncorrect user_id; "+err.Error())
		return
	}

	slugs, err := h.services.User.GetUserSegments(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, getUserSegmentsResponse{
		Status: "ok",
		Slugs:  slugs,
	})

}

type addUserToSegmentsRequest struct {
	UserID int            `json:"user_id" binding:"required"`
	Slugs  []string       `json:"slugs" binding:"required"`
	TTL    *time.Duration `json:"ttl"`
}

// @Summary Add User To Segment
// @Tags users
// @Description add user to chosen segments
// @ModuleID userAddUserToSegments
// @Accept  json
// @Produce  json
// @Param input body addUserToSegmentsRequest true "add to segments"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/segments [post]
func (h *HandlerV1) addUserToSegments(c *gin.Context) {
	const op = "http.v1.user.addUserToSegments"
	logger := h.logger.With(slog.String("op", op))

	var input addUserToSegmentsRequest

	if err := c.BindJSON(&input); err != nil {
		logger.Error("problem with bind json", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UserSegment.AddUserToSegments(input.UserID, input.Slugs, input.TTL); err != nil {
		// TODO: доработать обработку ошибок
		logger.Error("problem with add user to segment", sl.Err(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}

type removeUserFromSegmentsRequest struct {
	UserID int      `json:"user_id" binding:"required"`
	Slugs  []string `json:"slugs" binding:"required"`
}

// @Summary Remove User From Segments
// @Tags users
// @Description remove user from chosen segments
// @ModuleID userRemoveUserFromSegments
// @Accept  json
// @Produce  json
// @Param input body removeUserFromSegmentsRequest true "remove from segments"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/segments [delete]
func (h *HandlerV1) removeUserFromSegments(c *gin.Context) {
	const op = "http.v1.user.addUserToSegments"
	logger := h.logger.With(slog.String("op", op))

	var input removeUserFromSegmentsRequest

	if err := c.BindJSON(&input); err != nil {
		logger.Error("problem with bind json", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UserSegment.RemoveUserFromSegments(input.UserID, input.Slugs); err != nil {
		// TODO: доработать обработку ошибок
		logger.Error("problem with remove user from segment", sl.Err(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}

type getUserHistoryRequest struct {
	UserID    int    `json:"user_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required" time_format:"2006-01"`
	EndDate   string `json:"end_date" binding:"required" time_format:""`
}

type getUserHistoryResponse struct {
	URL string `json:"url"`
}

// @Summary Get User History
// @Tags segments
// @Description get user segments history
// @ModuleID getUserHistory
// @Accept  json
// @Produce  json
// @Param input body getUserHistoryRequest true "user history"
// @Success 200 {object} getUserHistoryResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/history [post]
func (h *HandlerV1) getUserHistory(c *gin.Context) {
	const op = "http.v1.user.getUserHistory"
	logger := h.logger.With(slog.String("op", op))

	var input getUserHistoryRequest

	if err := c.BindJSON(&input); err != nil {
		logger.Error("problem with bind json", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	startDate, err := time.Parse("2006-01", input.StartDate)
	if err != nil {
		logger.Error("problem uncorrect date", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	endDate, err := time.Parse("2006-01", input.EndDate)
	if err != nil {
		logger.Error("problem uncorrect date", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, err := h.services.UserSegment.GetUserSegmentHistory(input.UserID, startDate, endDate)
	if err != nil {
		logger.Error("problem with get user history", sl.Err(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getUserHistoryResponse{
		URL: h.basePath + "file/" + file,
	})

}
