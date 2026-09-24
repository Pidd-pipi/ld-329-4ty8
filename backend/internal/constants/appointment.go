package constants

// 预约状态：发起后等待对方确认，确认后进入待完成，交换结束标记完成，任一方可在完成前取消。
const (
	AppointmentStatusPending   = "PENDING"
	AppointmentStatusConfirmed = "CONFIRMED"
	AppointmentStatusCompleted = "COMPLETED"
	AppointmentStatusCancelled = "CANCELLED"
)

// 约见方式：线上会议需填写会议链接，线下需填写具体地点。
const (
	LocationTypeOnline  = "ONLINE"
	LocationTypeOffline = "OFFLINE"
)

const (
	// PairSeparator 预约双方展示名之间的分隔符。
	PairSeparator = " ↔ "
	// AgendaMaxLength 协商议程最大字数。
	AgendaMaxLength = 200
	// LocationMaxLength 会议链接或地点的最大长度。
	LocationMaxLength = 200
)

// AppointmentStatusLabels 状态机内部值与中文展示文案的映射。
var AppointmentStatusLabels = map[string]string{
	AppointmentStatusPending:   "待对方确认",
	AppointmentStatusConfirmed: "待完成",
	AppointmentStatusCompleted: "已完成",
	AppointmentStatusCancelled: "已取消",
}

// LocationTypeLabels 约见方式内部值与中文展示文案的映射。
var LocationTypeLabels = map[string]string{
	LocationTypeOnline:  "线上会议",
	LocationTypeOffline: "线下地点",
}

// IsActiveAppointmentStatus 判断状态是否仍占用时段。
// 待确认与待完成均属于生效预约，会阻止同一人在同一时段重复发起；已完成或已取消则释放时段。
func IsActiveAppointmentStatus(status string) bool {
	return status == AppointmentStatusPending || status == AppointmentStatusConfirmed
}

// 预约业务校验与冲突提示文案，集中维护，禁止散落在业务代码中。
const (
	MsgAppointmentNotFound   = "预约不存在或已被移除"
	MsgInvalidAppointmentID  = "预约编号无效"
	MsgSlotRequired          = "请选择一个双方共同时段"
	MsgSlotNotCommon         = "所选时段不在双方共同可用时间内"
	MsgLocationTypeInvalid   = "请选择线上会议或线下地点"
	MsgMeetingLinkInvalid    = "请填写有效的线上会议链接（需以 http:// 或 https:// 开头）"
	MsgMeetingLinkTooLong    = "会议链接长度不能超过 200 个字符"
	MsgPlaceRequired         = "请填写线下约见地点"
	MsgAgendaTooLong         = "协商议程不能超过 200 字"
	MsgConfirmOnlyPending    = "只有待对方确认的预约才能确认"
	MsgCompleteOnlyConfirmed = "只有已确认、待完成的预约才能标记完成"
	MsgCancelOnlyActive      = "预约已结束，不能再次取消"
	MsgMatchNotFound         = "匹配卡片不存在或已失效"
	MsgInvalidRequestBody    = "请求参数不完整或格式有误"
)

// MsgSlotConflict 时段冲突提示模板：被占用的人、时段、冲突对象、冲突预约状态。
// 明确告知原安排保留，引导改选其他共同时段。
const MsgSlotConflict = "时段冲突：%s 在「%s」已有进行中的预约（与 %s，%s）。原安排已保留，请改选其他共同时段。"
