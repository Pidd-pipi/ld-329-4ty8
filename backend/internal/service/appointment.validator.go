package service

import (
	"strings"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/model"
)

// validateCreateAppointment 校验发起预约请求：
// 共同时段必须来自匹配卡片，线上需填会议链接，线下需填具体地点，议程长度受限。
func validateCreateAppointment(request model.CreateAppointmentRequest, match model.Match) error {
	slot := strings.TrimSpace(request.Slot)
	if slot == "" {
		return errors.NewValidation(constants.MsgSlotRequired)
	}
	if !contains(match.CommonSlots, slot) {
		return errors.NewValidation(constants.MsgSlotNotCommon)
	}

	place := strings.TrimSpace(request.Place)
	switch request.LocationType {
	case constants.LocationTypeOnline:
		if len([]rune(place)) > constants.LocationMaxLength {
			return errors.NewValidation(constants.MsgMeetingLinkTooLong)
		}
		if !isMeetingLink(place) {
			return errors.NewValidation(constants.MsgMeetingLinkInvalid)
		}
	case constants.LocationTypeOffline:
		if place == "" {
			return errors.NewValidation(constants.MsgPlaceRequired)
		}
	default:
		return errors.NewValidation(constants.MsgLocationTypeInvalid)
	}

	if len([]rune(strings.TrimSpace(request.Agenda))) > constants.AgendaMaxLength {
		return errors.NewValidation(constants.MsgAgendaTooLong)
	}
	return nil
}

// contains 判断时段是否在匹配卡片的共同可用时间内。
func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

// isMeetingLink 校验线上会议链接，仅接受 http/https 协议地址。
func isMeetingLink(value string) bool {
	if value == "" {
		return false
	}
	return strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://")
}
