package controller

import (
	"net/http"
	"strconv"

	"BE-DroneMEQ-Violation-Management/domain-layer/service"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"
	"BE-DroneMEQ-Violation-Management/helpers/messages"
	"BE-DroneMEQ-Violation-Management/helpers/utils"

	"github.com/gin-gonic/gin"
)

type (
	ViolationMediaController interface {
		CreateVsMedia(ctx *gin.Context)
	}
	violationMediaController struct {
		violationMediaService service.ViolationMediaService
	}
)

func NewViolationMediaController(vs service.ViolationMediaService) ViolationMediaController {
	return &violationMediaController{
		violationMediaService: vs,
	}
}

func (c *violationMediaController) CreateVsMedia(ctx *gin.Context) {
	userData, exist := ctx.Get("userData")
	if !exist {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "user not found", nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	metadata, ok := userData.(external.MetadataDTO)
	if !ok {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "user not found", nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	violationIdStr := ctx.Param("Pkid")

	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_MEDIA, "failed to parse form data", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Extract multiple files from the form
	files := form.File["media"]
	if len(files) == 0 {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_MEDIA, "no media files found", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Create the violation media for all files
	results, err := c.violationMediaService.CreateViolationMedia(ctx.Request.Context(), files, violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_CREATE_VIOLATION_MEDIA, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_CREATE_VIOLATION_MEDIA, results)
	ctx.JSON(http.StatusOK, res)
}
