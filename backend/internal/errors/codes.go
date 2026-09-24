package errors

// 预约业务错误码，前端依据 code 做差异化提示（如 SLOT_CONFLICT 明确提示冲突）
const (
	CodeInvalidPayload    = "INVALID_PAYLOAD"
	CodeMatchNotFound     = "MATCH_NOT_FOUND"
	CodeActorNotParticip  = "ACTOR_NOT_PARTICIPANT"
	CodeInvalidCommonSlot = "INVALID_COMMON_SLOT"
	CodeSlotConflict      = "SLOT_CONFLICT"
	CodeApptNotFound      = "APPOINTMENT_NOT_FOUND"
	CodeNotConfirmer      = "NOT_CONFIRMER"
	CodeNotParticipant    = "NOT_PARTICIPANT"
	CodeInvalidStatus     = "INVALID_STATUS"
)

// 预约状态中文文案，供服务层拼装日志与返回信息
const (
	MsgInvalidPayload    = "预约信息不完整，请检查后重试"
	MsgMatchNotFound     = "匹配不存在，无法发起预约"
	MsgActorNotParticip  = "只有匹配双方可以发起该预约"
	MsgInvalidCommonSlot = "请选择双方的共同时段"
	MsgApptNotFound      = "预约不存在或已被删除"
	MsgNotConfirmer      = "只有被邀请的一方可以确认预约"
	MsgNotParticipant    = "只有预约双方可以执行该操作"
	MsgInvalidStatus     = "当前预约状态不允许该操作"
)

func SlotConflictMessage(owner, slot string) string {
	return "时段冲突：" + owner + " 在「" + slot + "」已有生效预约，原安排保持不变"
}
