package routes

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/controller"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group(constants.APIPrefix)
	api.GET("/health", controller.Health)
	api.GET("/dashboard/overview", controller.Overview)
	api.GET("/skills", controller.Skills)
	api.GET("/needs", controller.Needs)
	api.GET("/matches", controller.Matches)
	api.GET("/appointments", controller.Appointments)
	api.POST("/appointments", controller.CreateAppointment)
	api.POST("/appointments/:id/confirm", controller.ConfirmAppointment)
	api.POST("/appointments/:id/complete", controller.CompleteAppointment)
	api.POST("/appointments/:id/cancel", controller.CancelAppointment)
	api.GET("/reviews", controller.Reviews)
	api.GET("/messages", controller.Messages)
	api.GET("/profile", controller.Profile)
}
