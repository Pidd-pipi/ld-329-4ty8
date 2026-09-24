package controller

import (
	"net/http"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/logger"
	"cyskillswap/internal/model"
	"cyskillswap/internal/service"
	"github.com/gin-gonic/gin"
)

// ListAppointments GET /api/appointments：预约区展示发起人、确认人、状态与时段地点。
func ListAppointments(c *gin.Context) {
	c.JSON(http.StatusOK, service.ListAppointments())
}

// CreateAppointment POST /api/appointments：在匹配卡片中选择共同时段并发起预约。
func CreateAppointment(c *gin.Context) {
	var request model.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errors.NewValidation(constants.MsgInvalidRequestBody))
		return
	}
	appointment, err := service.CreateAppointment(request)
	if err != nil {
		respondError(c, err)
		return
	}
	logger.Info("appointment created", appointment.ID, appointment.Pair, appointment.Slot)
	c.JSON(http.StatusCreated, model.AppointmentMutationResult{
		Appointment: appointment, Appointments: service.ListAppointments(),
	})
}

// ConfirmAppointment POST /api/appointments/:id/confirm：对方确认，进入待完成。
func ConfirmAppointment(c *gin.Context) {
	mutateAppointment(c, service.ConfirmAppointment)
}

// CompleteAppointment POST /api/appointments/:id/complete：交换结束标记完成。
func CompleteAppointment(c *gin.Context) {
	mutateAppointment(c, service.CompleteAppointment)
}

// CancelAppointment POST /api/appointments/:id/cancel：任一方完成前取消，时段重新开放。
func CancelAppointment(c *gin.Context) {
	mutateAppointment(c, service.CancelAppointment)
}

// mutateAppointment 通用的预约状态流转入口，统一处理编号解析与响应封装。
func mutateAppointment(c *gin.Context, transition func(int) (model.Appointment, error)) {
	id, ok := parseAppointmentID(c)
	if !ok {
		respondError(c, errors.NewValidation(constants.MsgInvalidAppointmentID))
		return
	}
	appointment, err := transition(id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, model.AppointmentMutationResult{
		Appointment: appointment, Appointments: service.ListAppointments(),
	})
}
