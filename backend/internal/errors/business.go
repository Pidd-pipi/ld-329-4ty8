package errors

import "net/http"

// BusinessError 预约等业务流程的统一异常结构，错误码与 HTTP 状态配套返回给前端。
type BusinessError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e BusinessError) Error() string { return e.Message }

// 预约模块错误码，集中维护。
const (
	CodeValidation          = "VALIDATION_ERROR"
	CodeAppointmentNotFound = "APPOINTMENT_NOT_FOUND"
	CodeSlotConflict        = "SLOT_CONFLICT"
	CodeInvalidTransition   = "INVALID_STATUS_TRANSITION"
)

// NewValidation 参数校验失败。
func NewValidation(message string) BusinessError {
	return BusinessError{Status: http.StatusBadRequest, Code: CodeValidation, Message: message}
}

// NewConflict 同一参与者在同一时段已有生效预约，时段冲突。
func NewConflict(message string) BusinessError {
	return BusinessError{Status: http.StatusConflict, Code: CodeSlotConflict, Message: message}
}

// NewNotFound 预约记录不存在。
func NewNotFound(message string) BusinessError {
	return BusinessError{Status: http.StatusNotFound, Code: CodeAppointmentNotFound, Message: message}
}

// NewInvalidTransition 预约状态机不允许该流转。
func NewInvalidTransition(message string) BusinessError {
	return BusinessError{Status: http.StatusConflict, Code: CodeInvalidTransition, Message: message}
}
