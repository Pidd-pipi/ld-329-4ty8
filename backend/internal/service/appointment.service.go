package service

import (
	"fmt"
	"strings"
	"time"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/errors"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
)

const appointmentTimeLayout = "2006-01-02 15:04"

// nowText 统一生成预约记录使用的时间戳文本，便于测试时替换。
var nowText = func() string {
	return time.Now().Format(appointmentTimeLayout)
}

// ListAppointments 返回预约区展示所需的全量预约。
func ListAppointments() []model.Appointment {
	return repository.ListAppointments()
}

// CreateAppointment 在匹配卡片上发起预约：选择共同时段并填写线上会议链接或线下地点。
// 同一发起人或被约人在该时段已有生效预约时，明确返回冲突且保留原安排。
func CreateAppointment(request model.CreateAppointmentRequest) (model.Appointment, error) {
	match, ok := findMatch(request.MatchID)
	if !ok {
		return model.Appointment{}, errors.NewValidation(constants.MsgMatchNotFound)
	}
	if err := validateCreateAppointment(request, match); err != nil {
		return model.Appointment{}, err
	}

	slot := strings.TrimSpace(request.Slot)
	if err := ensureSlotFree(match.Provider, slot, 0); err != nil {
		return model.Appointment{}, err
	}
	if err := ensureSlotFree(match.Learner, slot, 0); err != nil {
		return model.Appointment{}, err
	}

	now := nowText()
	appointment := model.Appointment{
		MatchID:      match.ID,
		Initiator:    match.Learner,
		Confirmer:    match.Provider,
		Provider:     match.Provider,
		Learner:      match.Learner,
		OfferSkill:   match.OfferSkill,
		WantedSkill:  match.WantedSkill,
		Pair:         match.Provider + constants.PairSeparator + match.Learner,
		Slot:         slot,
		LocationType: request.LocationType,
		Place:        strings.TrimSpace(request.Place),
		Agenda:       strings.TrimSpace(request.Agenda),
		Status:       constants.AppointmentStatusPending,
		StatusLabel:  constants.AppointmentStatusLabels[constants.AppointmentStatusPending],
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return repository.SaveAppointment(appointment), nil
}

// ConfirmAppointment 被约人确认预约，预约由待确认进入待完成。
func ConfirmAppointment(id int) (model.Appointment, error) {
	return transitionAppointment(id, constants.AppointmentStatusConfirmed, constants.MsgConfirmOnlyPending)
}

// CompleteAppointment 交换结束后标记完成。
func CompleteAppointment(id int) (model.Appointment, error) {
	return transitionAppointment(id, constants.AppointmentStatusCompleted, constants.MsgCompleteOnlyConfirmed)
}

// CancelAppointment 任一方在完成前取消预约，取消后该时段可再次预约。
func CancelAppointment(id int) (model.Appointment, error) {
	return transitionAppointment(id, constants.AppointmentStatusCancelled, constants.MsgCancelOnlyActive)
}

// transitionAppointment 执行预约状态流转，前置状态不满足时返回业务异常。
func transitionAppointment(id int, target string, invalidMessage string) (model.Appointment, error) {
	appointment, ok := repository.GetAppointment(id)
	if !ok {
		return model.Appointment{}, errors.NewNotFound(constants.MsgAppointmentNotFound)
	}
	if !allowedTransition(appointment.Status, target) {
		return model.Appointment{}, errors.NewInvalidTransition(invalidMessage)
	}
	appointment.Status = target
	appointment.StatusLabel = constants.AppointmentStatusLabels[target]
	appointment.UpdatedAt = nowText()
	if !repository.UpdateAppointment(appointment) {
		return model.Appointment{}, errors.NewNotFound(constants.MsgAppointmentNotFound)
	}
	return appointment, nil
}

// allowedTransition 预约状态机：
// 待确认 -> 待完成 / 已取消；待完成 -> 已完成 / 已取消；已完成、已取消为终态。
func allowedTransition(from, to string) bool {
	switch from {
	case constants.AppointmentStatusPending:
		return to == constants.AppointmentStatusConfirmed || to == constants.AppointmentStatusCancelled
	case constants.AppointmentStatusConfirmed:
		return to == constants.AppointmentStatusCompleted || to == constants.AppointmentStatusCancelled
	default:
		return false
	}
}

// ensureSlotFree 检查参与者在该时段是否已被生效预约占用；若占用则返回冲突提示，原安排保留不变。
func ensureSlotFree(participant, slot string, excludeID int) error {
	conflict, occupied := repository.ActiveAppointmentFor(participant, slot, excludeID)
	if !occupied {
		return nil
	}
	other := conflict.Initiator
	if other == participant {
		other = conflict.Confirmer
	}
	return errors.NewConflict(fmt.Sprintf(
		constants.MsgSlotConflict,
		participant, slot, other, constants.AppointmentStatusLabels[conflict.Status],
	))
}

// findMatch 按编号查找智能匹配卡片。
func findMatch(id int) (model.Match, bool) {
	for _, match := range repository.ListMatches() {
		if match.ID == id {
			return match, true
		}
	}
	return model.Match{}, false
}
