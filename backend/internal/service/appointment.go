package service

import (
	"strings"
	"time"

	"cyskillswap/internal/constants"
	apperrors "cyskillswap/internal/errors"
	"cyskillswap/internal/logger"
	"cyskillswap/internal/model"
	"cyskillswap/internal/repository"
	"cyskillswap/internal/validator"
)

const appointmentTimeLayout = "2006-01-02 15:04"

// AppointmentError 预约业务错误，携带错误码供控制器映射 HTTP 状态
type AppointmentError struct {
	Code string
	Err  error
}

func (e AppointmentError) Error() string { return e.Err.Error() }

// Unwrap 暴露内部业务错误，便于控制器 errors.As 提取 code/message
func (e AppointmentError) Unwrap() error { return e.Err }

func conflictError(code, message string) AppointmentError {
	return AppointmentError{Code: code, Err: apperrors.BusinessError{Code: code, Message: message}}
}

// ListAppointments 返回预约时间线
func ListAppointments() []model.Appointment {
	return repository.ListAppointments()
}

// CreateAppointment 发起预约：
// 仅匹配双方可发起；时段必须是共同时段；任一方已有生效预约占用该时段时拒绝并保留原安排。
func CreateAppointment(req model.CreateAppointmentRequest) (model.Appointment, error) {
	if err := validator.ValidateCreateRequest(req); err != nil {
		return model.Appointment{}, AppointmentError{Code: apperrors.CodeInvalidPayload, Err: err}
	}

	match, ok := repository.FindMatch(req.MatchID)
	if !ok {
		return model.Appointment{}, conflictError(apperrors.CodeMatchNotFound, apperrors.MsgMatchNotFound)
	}
	if !validator.IsParticipant(match.Provider, match.Learner, req.Initiator) {
		return model.Appointment{}, conflictError(apperrors.CodeActorNotParticip, apperrors.MsgActorNotParticip)
	}
	if !validator.IsCommonSlot(match, req.Slot) {
		return model.Appointment{}, conflictError(apperrors.CodeInvalidCommonSlot, apperrors.MsgInvalidCommonSlot)
	}

	otherParty := match.Provider
	if req.Initiator == match.Provider {
		otherParty = match.Learner
	}
	// 任一方在该时段已有生效预约都要明确提示冲突，原安排保持不变
	for _, person := range []string{req.Initiator, otherParty} {
		if existing, conflict := repository.FindActiveConflict(person, req.Slot); conflict {
			owner := existing.Provider
			if existing.Learner == person {
				owner = existing.Learner
			}
			message := apperrors.SlotConflictMessage(owner, req.Slot)
			logger.Warn("appointment slot conflict", "initiator", req.Initiator, "slot", req.Slot, "keptAppt", existing.ID)
			return model.Appointment{}, conflictError(apperrors.CodeSlotConflict, message)
		}
	}

	appointment := model.Appointment{
		MatchID:     match.ID,
		Initiator:   req.Initiator,
		Provider:    match.Provider,
		Learner:     match.Learner,
		Pair:        match.Provider + " ↔ " + match.Learner,
		OfferSkill:  match.OfferSkill,
		WantedSkill: match.WantedSkill,
		Slot:        req.Slot,
		MeetingMode: req.MeetingMode,
		Location:    strings.TrimSpace(req.Location),
		Agenda:      strings.TrimSpace(req.Agenda),
		Status:      constants.AppointmentStatusPending,
		CreatedAt:   time.Now().Format(appointmentTimeLayout),
	}
	created := repository.CreateAppointment(appointment)
	logger.Info("appointment created", "id", created.ID, "initiator", req.Initiator, "slot", req.Slot)
	return created, nil
}

// ConfirmAppointment 被邀请方确认预约，状态由待确认进入待完成（confirmed）
func ConfirmAppointment(id int64, req model.AppointmentActionRequest) (model.Appointment, error) {
	if err := validator.ValidateActionRequest(req); err != nil {
		return model.Appointment{}, AppointmentError{Code: apperrors.CodeInvalidPayload, Err: err}
	}
	appointment, ok := repository.FindAppointment(id)
	if !ok {
		return model.Appointment{}, conflictError(apperrors.CodeApptNotFound, apperrors.MsgApptNotFound)
	}
	// 只有发起人之外的另一方可以确认
	if !validator.IsParticipant(appointment.Provider, appointment.Learner, req.Actor) {
		return model.Appointment{}, conflictError(apperrors.CodeNotParticipant, apperrors.MsgNotParticipant)
	}
	if req.Actor == appointment.Initiator {
		return model.Appointment{}, conflictError(apperrors.CodeNotConfirmer, apperrors.MsgNotConfirmer)
	}
	if appointment.Status != constants.AppointmentStatusPending {
		return model.Appointment{}, conflictError(apperrors.CodeInvalidStatus, apperrors.MsgInvalidStatus)
	}

	appointment.Confirmer = req.Actor
	appointment.Status = constants.AppointmentStatusConfirmed
	appointment.ConfirmedAt = time.Now().Format(appointmentTimeLayout)
	repository.UpdateAppointment(appointment)
	logger.Info("appointment confirmed", "id", id, "confirmer", req.Actor)
	return appointment, nil
}

// CompleteAppointment 交换结束后由任一方标记完成
func CompleteAppointment(id int64, req model.AppointmentActionRequest) (model.Appointment, error) {
	if err := validator.ValidateActionRequest(req); err != nil {
		return model.Appointment{}, AppointmentError{Code: apperrors.CodeInvalidPayload, Err: err}
	}
	appointment, ok := repository.FindAppointment(id)
	if !ok {
		return model.Appointment{}, conflictError(apperrors.CodeApptNotFound, apperrors.MsgApptNotFound)
	}
	if !validator.IsParticipant(appointment.Provider, appointment.Learner, req.Actor) {
		return model.Appointment{}, conflictError(apperrors.CodeNotParticipant, apperrors.MsgNotParticipant)
	}
	if appointment.Status != constants.AppointmentStatusConfirmed {
		return model.Appointment{}, conflictError(apperrors.CodeInvalidStatus, apperrors.MsgInvalidStatus)
	}

	appointment.Status = constants.AppointmentStatusCompleted
	appointment.CompletedAt = time.Now().Format(appointmentTimeLayout)
	repository.UpdateAppointment(appointment)
	logger.Info("appointment completed", "id", id, "actor", req.Actor)
	return appointment, nil
}

// CancelAppointment 完成前任一方可取消；取消后该时段释放，可再次发起预约
func CancelAppointment(id int64, req model.AppointmentActionRequest) (model.Appointment, error) {
	if err := validator.ValidateActionRequest(req); err != nil {
		return model.Appointment{}, AppointmentError{Code: apperrors.CodeInvalidPayload, Err: err}
	}
	appointment, ok := repository.FindAppointment(id)
	if !ok {
		return model.Appointment{}, conflictError(apperrors.CodeApptNotFound, apperrors.MsgApptNotFound)
	}
	if !validator.IsParticipant(appointment.Provider, appointment.Learner, req.Actor) {
		return model.Appointment{}, conflictError(apperrors.CodeNotParticipant, apperrors.MsgNotParticipant)
	}
	if appointment.Status != constants.AppointmentStatusPending && appointment.Status != constants.AppointmentStatusConfirmed {
		return model.Appointment{}, conflictError(apperrors.CodeInvalidStatus, apperrors.MsgInvalidStatus)
	}

	appointment.Status = constants.AppointmentStatusCancelled
	appointment.CancelledBy = req.Actor
	appointment.CancelledAt = time.Now().Format(appointmentTimeLayout)
	repository.UpdateAppointment(appointment)
	logger.Info("appointment cancelled", "id", id, "actor", req.Actor, "slot released", appointment.Slot)
	return appointment, nil
}

// ActiveAppointmentCount 统计生效中（待确认 + 待完成）的预约数量，用于首页指标
func ActiveAppointmentCount(appointments []model.Appointment) int {
	count := 0
	for _, appointment := range appointments {
		if constants.IsActiveAppointmentStatus(appointment.Status) {
			count++
		}
	}
	return count
}
