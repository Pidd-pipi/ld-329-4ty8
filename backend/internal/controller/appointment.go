package controller

import (
	stderrors "errors"
	"net/http"
	"strconv"

	apperrors "cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

// writeAppointmentError 把预约业务错误码映射为合适的 HTTP 状态，
// 前端可从返回体的 code/message 中给出明确提示（如时段冲突）。
func writeAppointmentError(c *gin.Context, err error) {
	var appointmentErr service.AppointmentError
	if !stderrors.As(err, &appointmentErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "服务暂时不可用"})
		return
	}

	status := http.StatusBadRequest
	switch appointmentErr.Code {
	case apperrors.CodeMatchNotFound, apperrors.CodeApptNotFound:
		status = http.StatusNotFound
	case apperrors.CodeSlotConflict:
		status = http.StatusConflict
	case apperrors.CodeActorNotParticip, apperrors.CodeNotConfirmer, apperrors.CodeNotParticipant:
		status = http.StatusForbidden
	}

	var businessErr apperrors.BusinessError
	if stderrors.As(err, &businessErr) {
		c.JSON(status, gin.H{"code": businessErr.Code, "message": businessErr.Message})
		return
	}
	c.JSON(status, gin.H{"code": appointmentErr.Code, "message": appointmentErr.Error()})
}

// CreateAppointment POST /api/appointments
func CreateAppointment(c *gin.Context) {
	var req model.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": apperrors.CodeInvalidPayload, "message": apperrors.MsgInvalidPayload})
		return
	}
	appointment, err := service.CreateAppointment(req)
	if err != nil {
		writeAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusCreated, appointment)
}

// ConfirmAppointment POST /api/appointments/:id/confirm
func ConfirmAppointment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": apperrors.CodeInvalidPayload, "message": apperrors.MsgApptNotFound})
		return
	}
	var req model.AppointmentActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": apperrors.CodeInvalidPayload, "message": apperrors.MsgInvalidPayload})
		return
	}
	appointment, err := service.ConfirmAppointment(id, req)
	if err != nil {
		writeAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, appointment)
}

// CompleteAppointment POST /api/appointments/:id/complete
func CompleteAppointment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": apperrors.CodeInvalidPayload, "message": apperrors.MsgApptNotFound})
		return
	}
	var req model.AppointmentActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": apperrors.CodeInvalidPayload, "message": apperrors.MsgInvalidPayload})
		return
	}
	appointment, err := service.CompleteAppointment(id, req)
	if err != nil {
		writeAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, appointment)
}

// CancelAppointment POST /api/appointments/:id/cancel
func CancelAppointment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": apperrors.CodeInvalidPayload, "message": apperrors.MsgApptNotFound})
		return
	}
	var req model.AppointmentActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": apperrors.CodeInvalidPayload, "message": apperrors.MsgInvalidPayload})
		return
	}
	appointment, err := service.CancelAppointment(id, req)
	if err != nil {
		writeAppointmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, appointment)
}
