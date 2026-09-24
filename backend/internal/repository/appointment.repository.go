package repository

import (
	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

// seed 初始化两条预约演示数据：一条已确认待完成，一条待对方确认。
// 参与者与 ListMatches 保持一致，用于验证同一时段冲突检测。
func (s *appointmentStore) seed() {
	seeds := []model.Appointment{
		{
			ID: 1, MatchID: 1, Initiator: "孟野", Confirmer: "林澈",
			Provider: "林澈", Learner: "孟野", OfferSkill: "毕业照人像摄影", WantedSkill: "民谣吉他陪练",
			Pair: "林澈 ↔ 孟野", Slot: "周六上午", LocationType: constants.LocationTypeOffline,
			Place: "东校区湖边", Agenda: "先拍宣传照，再约 2 次吉他课",
			Status: constants.AppointmentStatusConfirmed, StatusLabel: constants.AppointmentStatusLabels[constants.AppointmentStatusConfirmed],
			CreatedAt: "2026-09-20 09:30", UpdatedAt: "2026-09-21 12:10",
		},
		{
			ID: 2, MatchID: 2, Initiator: "许安", Confirmer: "周芮",
			Provider: "周芮", Learner: "许安", OfferSkill: "Python 数据分析", WantedSkill: "论文数据清洗",
			Pair: "周芮 ↔ 许安", Slot: "周二晚", LocationType: constants.LocationTypeOnline,
			Place: "https://meeting.example.com/room/python-2048", Agenda: "导入问卷 CSV 并完成基础可视化",
			Status: constants.AppointmentStatusPending, StatusLabel: constants.AppointmentStatusLabels[constants.AppointmentStatusPending],
			CreatedAt: "2026-09-22 20:15", UpdatedAt: "2026-09-22 20:15",
		},
	}
	s.appointments = seeds
	s.sequence = len(seeds)
}

// ListAppointments 返回全部预约（按编号倒序，最新发起的在前）。
func ListAppointments() []model.Appointment {
	result := appointmentData.read(func() any {
		items := appointmentData.snapshot()
		return items
	}).([]model.Appointment)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// GetAppointment 按编号查询预约。
func GetAppointment(id int) (model.Appointment, bool) {
	var found *model.Appointment
	appointmentData.read(func() any {
		if got, ok := appointmentData.findByID(id); ok {
			found = &got
		}
		return nil
	})
	if found == nil {
		return model.Appointment{}, false
	}
	return *found, true
}

// SaveAppointment 新增预约，落库后返回带编号的完整记录。
func SaveAppointment(appointment model.Appointment) model.Appointment {
	appointmentData.write(func() {
		appointment.ID = appointmentData.nextID()
		appointmentData.add(appointment)
	})
	return appointment
}

// UpdateAppointment 按编号更新预约状态等信息。
func UpdateAppointment(appointment model.Appointment) bool {
	ok := false
	appointmentData.write(func() {
		ok = appointmentData.replace(appointment)
	})
	return ok
}

// ActiveAppointmentFor 检查某位参与者在指定时段是否已有生效预约。
// 生效预约指待确认或待完成；excludeID 用于忽略自身（取消后状态变更不影响检测）。
// 返回冲突预约及其是否存在。
func ActiveAppointmentFor(participant, slot string, excludeID int) (model.Appointment, bool) {
	var conflict *model.Appointment
	appointmentData.read(func() any {
		for _, appointment := range appointmentData.appointments {
			if appointment.ID == excludeID {
				continue
			}
			if appointment.Slot != slot || !constants.IsActiveAppointmentStatus(appointment.Status) {
				continue
			}
			if appointment.Initiator == participant || appointment.Confirmer == participant {
				got := appointment
				conflict = &got
				return nil
			}
		}
		return nil
	})
	if conflict == nil {
		return model.Appointment{}, false
	}
	return *conflict, true
}
