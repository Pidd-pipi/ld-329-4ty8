package controller

import (
	"net/http"
	"strconv"

	"cyskillswap/internal/errors"
	"cyskillswap/internal/logger"
	"github.com/gin-gonic/gin"
)

// parseAppointmentID 从路径参数解析预约编号。
func parseAppointmentID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// respondError 将业务异常翻译成统一的错误响应，未知错误按 500 兜底。
func respondError(c *gin.Context, err error) {
	if business, ok := err.(errors.BusinessError); ok {
		logger.Warn("appointment business error", business.Code, business.Message)
		c.JSON(business.Status, business)
		return
	}
	logger.Error("appointment unexpected error", err)
	c.JSON(http.StatusInternalServerError, errors.BusinessError{
		Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "服务暂时不可用",
	})
}
