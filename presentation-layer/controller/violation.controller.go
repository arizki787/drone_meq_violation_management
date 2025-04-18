package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"BE-DroneMEQ-Violation-Management/domain-layer/service"
	"BE-DroneMEQ-Violation-Management/helpers/dto"
	"BE-DroneMEQ-Violation-Management/helpers/dto/external"
	"BE-DroneMEQ-Violation-Management/helpers/messages"
	"BE-DroneMEQ-Violation-Management/helpers/utils"

	"github.com/gin-gonic/gin"
)

type (
	ViolationController interface {
		CreateVs(ctx *gin.Context)
		GetAllViolation(ctx *gin.Context)
		ViolationById(ctx *gin.Context)
		ViolationByFilter(ctx *gin.Context)
		UpdateVs(ctx *gin.Context)
		SoftDeleteVs(ctx *gin.Context)
		RestoreVs(ctx *gin.Context)
		DeleteVs(ctx *gin.Context)
	}

	violationController struct {
		violationService service.ViolationService
	}
)

func NewViolationController(vs service.ViolationService) ViolationController {
	return &violationController{
		violationService: vs,
	}
}

func (c *violationController) CreateVs(ctx *gin.Context) {
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

	var req dto.ViolationCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	violationResponse, err := c.violationService.CreateViolation(ctx.Request.Context(), req, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_CREATE_VIOLATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	ctx.JSON(http.StatusCreated, utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_CREATE_VIOLATION, violationResponse))
}

func (c *violationController) GetAllViolation(ctx *gin.Context) {
	userData, exist := ctx.Get("userData")
	if !exist {
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
	var req dto.PaginationRequest
	// Check if the request body is valid
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Get all violations with pagination
	result, err := c.violationService.GetAllViolationWithPagination(ctx.Request.Context(), req, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_ALL_VIOLATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_GET_ALL_VIOLATION, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationController) ViolationById(ctx *gin.Context) {
	// Get Pkid from URL parameter
	violationIdStr := ctx.Param("Pkid")

	// Convert Pkid to int64
	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		// Handle the case where Pkid is not a valid integer
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Fetch the violation using the violationService
	result, err := c.violationService.GetViolationById(ctx.Request.Context(), violationId)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_VIOLATION_BY_ID, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Return the result
	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_GET_VIOLATION_BY_ID, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationController) ViolationByFilter(ctx *gin.Context) {
	userData, exist := ctx.Get("userData")
	if !exist {
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
	// Get status from URL parameter
	status := ctx.QueryArray("status")
	violationType := ctx.QueryArray("violation_type")
	severity := ctx.QueryArray("severity")

	var req dto.PaginationRequest
	// Check if the request body is valid
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil))
		return
	}

	result, err := c.violationService.GetViolationByFilter(ctx.Request.Context(), status, violationType, severity, req, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_VIOLATION_BY_STATUS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
	}

	// Return the result
	ctx.JSON(http.StatusOK, dto.ViolationPaginationResponse{
		Data:               result.Data,
		PaginationResponse: result.PaginationResponse,
	})
}

func (c *violationController) UpdateVs(ctx *gin.Context) {
	userData, exist := ctx.Get("userData")
	if !exist {
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

	// Get Pkid from URL parameter
	violationIdStr := ctx.Param("Pkid")

	// Convert Pkid to int64
	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		// Handle the case where Pkid is not a valid integer
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Get the request body
	var req dto.ViolationUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Call the violationService to update the violation
	result, err := c.violationService.UpdateViolation(ctx.Request.Context(), req, violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_UPDATE_VIOLATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Return the result
	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_UPDATE_VIOLATION, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationController) SoftDeleteVs(ctx *gin.Context) {
	userData, exist := ctx.Get("userData")
	// Check if the user data exists
	if !exist {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "User data not found", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	metadata, ok := userData.(external.MetadataDTO)
	// Check if the user data is valid
	if !ok {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "User data not valid", nil)
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

	result, err := c.violationService.SoftDeleteViolation(ctx.Request.Context(), violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_SOFT_DELETE_VIOLATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_SOFT_DELETE_VIOLATION, result))
}

func (c *violationController) RestoreVs(ctx *gin.Context) {
	userData, exist := ctx.Get("userData")
	if !exist {
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

	result, err := c.violationService.RestoreViolation(ctx.Request.Context(), violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_RESTORE_VIOLATION, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_RESTORE_VIOLATION, result))
}

func (c *violationController) DeleteVs(ctx *gin.Context) {
	violationIdStr := ctx.Param("Pkid")
	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.violationService.HardDeleteViolation(ctx.Request.Context(), violationId); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusNotFound, utils.BuildResponseFailed(messages.MESSAGE_FAILED_DELETE_VIOLATION, err.Error(), nil))
		} else {
			log.Printf("Failed to delete violation with Pkid: %v", violationId)
			ctx.JSON(http.StatusInternalServerError, utils.BuildResponseFailed(messages.MESSAGE_FAILED_DELETE_VIOLATION, err.Error(), nil))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_DELETE_VIOLATION, nil))
}
