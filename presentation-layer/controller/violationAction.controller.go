package controller

import (
	"net/http"
	"strconv"

	"BE-DroneMEQ-Violation-Management/domain-layer/service"
	"BE-DroneMEQ-Violation-Management/helpers/dto"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"
	"BE-DroneMEQ-Violation-Management/helpers/messages"
	"BE-DroneMEQ-Violation-Management/helpers/utils"

	"github.com/gin-gonic/gin"
)

type (
	ViolationActionController interface {
		CreateVsAction(ctx *gin.Context)
	}

	violationActionController struct {
		violationActionService service.ViolationActionService
	}
)

func NewViolationActionController(vs service.ViolationActionService) ViolationActionController {
	return &violationActionController{
		violationActionService: vs,
	}
}

func (c *violationActionController) CreateVsAction(ctx *gin.Context) {
	userData, exists := ctx.Get("userData")
	if !exists {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "User data not found", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	metadata, ok := userData.(external.MetadataDTO)
	if !ok {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "User data not found", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	violationIdStr := ctx.Param("Pkid")

	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var violationAction dto.ViolationActionCreateRequest
	// Check if the request body is valid
	if err := ctx.ShouldBind(&violationAction); err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Create the violation action
	result, err := c.violationActionService.CreateViolationAction(ctx.Request.Context(), violationAction, violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_CREATE_VIOLATION_ACTION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	//return the result
	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_CREATE_VIOLATION_ACTION, result)
	ctx.JSON(http.StatusOK, res)
}
