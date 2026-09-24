package constants

// 预约状态：贯穿「待确认 → 待完成 → 已完成」主流程，取消为终止态
const (
	AppointmentStatusPending   = "pending"
	AppointmentStatusConfirmed = "confirmed"
	AppointmentStatusCompleted = "completed"
	AppointmentStatusCancelled = "cancelled"
)

// 预约形式：线上会议或线下面对面
const (
	MeetingModeOnline  = "online"
	MeetingModeOffline = "offline"
)

// 预约操作类型（确认 / 完成 / 取消）
const (
	AppointmentActionConfirm   = "confirm"
	AppointmentActionComplete  = "complete"
	AppointmentActionCancel    = "cancel"
	AppointmentActionActorSelf = "self"
)

// 占用时段、阻断新申请的「生效中」状态
var ActiveAppointmentStatuses = []string{
	AppointmentStatusPending,
	AppointmentStatusConfirmed,
}

func IsActiveAppointmentStatus(status string) bool {
	for _, active := range ActiveAppointmentStatuses {
		if status == active {
			return true
		}
	}
	return false
}
