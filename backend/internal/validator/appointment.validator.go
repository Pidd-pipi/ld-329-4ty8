package validator

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/model"
)

// ValidateCreateRequest 校验发起预约入参本身：形式合法且地点/链接非空
func ValidateCreateRequest(req model.CreateAppointmentRequest) error {
	if req.MeetingMode != constants.MeetingModeOnline && req.MeetingMode != constants.MeetingModeOffline {
		return errors.BusinessError{Code: errors.CodeInvalidPayload, Message: errors.MsgInvalidPayload}
	}
	if req.Location == "" {
		return errors.BusinessError{Code: errors.CodeInvalidPayload, Message: "请填写线上会议链接或线下地点"}
	}
	return nil
}

// ValidateActionRequest 校验确认/完成/取消入参
func ValidateActionRequest(req model.AppointmentActionRequest) error {
	if req.Actor == "" {
		return errors.BusinessError{Code: errors.CodeInvalidPayload, Message: errors.MsgInvalidPayload}
	}
	return nil
}

// IsCommonSlot 判断所选时段是否属于匹配双方的共同时段
func IsCommonSlot(match model.Match, slot string) bool {
	for _, common := range match.CommonSlots {
		if common == slot {
			return true
		}
	}
	return false
}

// IsParticipant 判断某人是否为匹配/预约的当事人
func IsParticipant(provider, learner, person string) bool {
	return person == provider || person == learner
}
