package model

// CreateAppointmentRequest 发起预约入参：从匹配卡片选择共同时段并填写会议链接/线下地点
type CreateAppointmentRequest struct {
	MatchID     int    `json:"matchId"     binding:"required"`
	Initiator   string `json:"initiator"   binding:"required"`
	Slot        string `json:"slot"        binding:"required"`
	MeetingMode string `json:"meetingMode" binding:"required"`
	Location    string `json:"location"    binding:"required"`
	Agenda      string `json:"agenda"`
}

// AppointmentActionRequest 确认 / 完成 / 取消入参，携带当前操作人用于权限校验
type AppointmentActionRequest struct {
	Actor string `json:"actor" binding:"required"`
}
