package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/user-segmentation-service/pkg/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *HandlerV1) initSegmentRoute(api *gin.RouterGroup) {
	segment := api.Group("/segment")
	{
		segment.GET("", h.getSegment)
		segment.POST("", h.createSegment)
	}
}

// @Summary Get Segment
// @Tags segments
// @Description get segment info
// @ModuleID getSegment
// @Accept  json
// @Produce  json
//
//	Param slug path string true "segment slug"
//
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments/{slug} [get]
func (h *HandlerV1) getSegment(c *gin.Context) {

}

type createSegmentRequest struct {
	Slug             string `json:"slug" binding:"required,slug"`
	AssignPercentage int    `json:"assign_percentage" binding:"gte=0,lte=100"`
}

// @Summary Create Segment
// @Tags segments
// @Description create new segment
// @ModuleID createSegment
// @Accept  json
// @Produce  json
// @Param input body createSegmentRequest true "create segment"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments [post]
func (h *HandlerV1) createSegment(c *gin.Context) {
	const op = "http.v1.segment.createSegment"
	logger := h.logger.With(slog.String("op", op))

	var input createSegmentRequest

	if err := c.BindJSON(&input); err != nil {
		logger.Error("problem with bind json", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	Id, err := h.services.Segment.CreateSegment(input.Slug)
	if err != nil {
		logger.Error("problem createElement", sl.Err(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if input.AssignPercentage > 0 {
		if err := h.services.UserSegment.SetSegmentToRandomUsers(input.Slug, input.AssignPercentage); err != nil {
			logger.Error("problem SetSegmentToRandomUsers", sl.Err(err))
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"id":     strconv.Itoa(Id),
	})
}

type deleteSegmentInput struct {
	Slug string `json:"slug" binding:"required"`
}

// @Summary Delete Segment
// @Tags segments
// @Description delete segment
// @ModuleID deleteSegment
// @Accept  json
// @Produce  json
// @Param input body deleteSegmentInput true "delete segment"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /segments [delete]
func (h *HandlerV1) deleteSegment(c *gin.Context) {
	const op = "http.v1.segment.deleteSegment"
	logger := h.logger.With(slog.String("op", op))

	var input deleteSegmentInput

	if err := c.BindJSON(&input); err != nil {
		logger.Error("problem with bind json", sl.Err(err))
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Segment.DeleteSegment(input.Slug); err != nil {
		logger.Error("problem createElement", sl.Err(err))
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
