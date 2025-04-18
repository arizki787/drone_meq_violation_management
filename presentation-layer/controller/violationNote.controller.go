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
	ViolationNoteController interface {
		CreateVsNote(ctx *gin.Context)
		GetVsNote(ctx *gin.Context)
		GetVsNoteByFilter(ctx *gin.Context)
		UpdateVsNote(ctx *gin.Context)
		SoftDeleteVsNote(ctx *gin.Context)
		RestoreVsNote(ctx *gin.Context)
		DeleteVsNote(ctx *gin.Context)
	}
	violationNoteController struct {
		violationNoteService service.ViolationNoteService
	}
)

func NewViolationNoteController(vs service.ViolationNoteService) ViolationNoteController {
	return &violationNoteController{
		violationNoteService: vs,
	}
}

func (c *violationNoteController) CreateVsNote(ctx *gin.Context) {
	violationIdStr := ctx.Param("Pkid")
	userData, exist := ctx.Get("userData")
	if !exist {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "User data not found", nil)
		ctx.JSON(http.StatusBadRequest, res)
	}
	metadata, ok := userData.(external.MetadataDTO)
	if !ok {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_USER_DATA, "User data not found", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var violationNote dto.ViolationNoteCreateRequest
	// Check if the request body is valid
	if err := ctx.ShouldBind(&violationNote); err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Create the violation note
	result, err := c.violationNoteService.CreateViolationNote(ctx.Request.Context(), violationNote, violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_CREATE_VIOLATION_NOTE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	//return the result
	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_CREATE_VIOLATION_NOTE, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationNoteController) GetVsNote(ctx *gin.Context) {
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

	result, err := c.violationNoteService.GetViolationNoteByViolationId(ctx.Request.Context(), dto.PaginationRequest{}, violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_VIOLATION_NOTE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_GET_VIOLATION_NOTE, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationNoteController) GetVsNoteByFilter(ctx *gin.Context) {
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

	noteType := ctx.QueryArray("note_type")

	var req dto.PaginationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil))
		return
	}

	result, err := c.violationNoteService.GetViolationNoteByFilter(ctx.Request.Context(), noteType, dto.PaginationRequest{}, violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_VIOLATION_NOTE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_GET_VIOLATION_NOTE, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationNoteController) UpdateVsNote(ctx *gin.Context) {
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
	violationNoteIdStr := ctx.Param("notePkid")

	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	violationNoteId, err := strconv.ParseInt(violationNoteIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var violationNote dto.ViolationNoteUpdateRequest
	// Check if the request body is valid
	if err := ctx.ShouldBind(&violationNote); err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Update the violation note
	result, err := c.violationNoteService.UpdateViolationNote(ctx.Request.Context(), violationNote, violationNoteId, violationId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_UPDATE_VIOLATION_NOTE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	//return the result
	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_UPDATE_VIOLATION_NOTE, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationNoteController) SoftDeleteVsNote(ctx *gin.Context) {
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
	violationNoteIdStr := ctx.Param("notePkid")

	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	violationNoteId, err := strconv.ParseInt(violationNoteIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Soft delete the violation note
	result, err := c.violationNoteService.SoftDeleteViolationNote(ctx.Request.Context(), violationId, violationNoteId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_SOFT_DELETE_VIOLATION_NOTE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	//return the result
	res := utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_SOFT_DELETE_VIOLATION_NOTE, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *violationNoteController) RestoreVsNote(ctx *gin.Context) {
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
	violationNoteIdStr := ctx.Param("notePkid")

	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	violationNoteId, err := strconv.ParseInt(violationNoteIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.violationNoteService.RestoreViolationNote(ctx.Request.Context(), violationId, violationNoteId, metadata)
	if err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_RESTORE_VIOLATION_NOTE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_RESTORE_VIOLATION_NOTE, result))

}

func (c *violationNoteController) DeleteVsNote(ctx *gin.Context) {
	violationIdStr := ctx.Param("Pkid")
	violationNoteIdStr := ctx.Param("notePkid")

	violationId, err := strconv.ParseInt(violationIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	violationNoteId, err := strconv.ParseInt(violationNoteIdStr, 10, 64)
	if err != nil {
		res := utils.BuildResponseFailed("Invalid Pkid", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.violationNoteService.HardDeleteViolationNote(ctx.Request.Context(), violationId, violationNoteId); err != nil {
		res := utils.BuildResponseFailed(messages.MESSAGE_FAILED_HARD_DELETE_VIOLATION_NOTE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	ctx.JSON(http.StatusOK, utils.BuildResponseSuccess(messages.MESSAGE_SUCCESS_HARD_DELETE_VIOLATION_NOTE, nil))
}
