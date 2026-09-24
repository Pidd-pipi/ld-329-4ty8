package repository

import (
	"sync"

	"cyskillswap/internal/constants"
	"cyskillswap/internal/model"
)

// appointmentStore 预约的内存持久层。
// 单进程内通过互斥锁保证并发安全；状态为 pending/confirmed 的预约占用时段。
type appointmentStore struct {
	mu           sync.RWMutex
	nextID       int64
	appointments []model.Appointment
}

var store = newAppointmentStore()

func newAppointmentStore() *appointmentStore {
	return &appointmentStore{
		nextID: 3,
		appointments: []model.Appointment{
			{
				ID: 1, MatchID: 3, Initiator: "孟野", Confirmer: "林澈",
				Provider: "孟野", Learner: "林澈", Pair: "孟野 ↔ 林澈",
				OfferSkill: "民谣吉他陪练", WantedSkill: "宣传照拍摄",
				Slot: "周六上午", MeetingMode: constants.MeetingModeOffline, Location: "东校区湖边",
				Agenda: "先拍宣传照，再约 2 次吉他课",
				Status: constants.AppointmentStatusConfirmed,
				CreatedAt: "2026-09-20 09:30", ConfirmedAt: "2026-09-20 21:10",
			},
			{
				ID: 2, MatchID: 2, Initiator: "许安",
				Provider: "周芮", Learner: "许安", Pair: "周芮 ↔ 许安",
				OfferSkill: "Python 数据分析", WantedSkill: "论文数据清洗",
				Slot: "周二晚", MeetingMode: constants.MeetingModeOnline, Location: "https://meeting.example.com/room/python-pandas",
				Agenda: "导入问卷 CSV 并完成基础可视化",
				Status: constants.AppointmentStatusPending,
				CreatedAt: "2026-09-22 18:05",
			},
		},
	}
}

// ListAppointments 返回全部预约（最新发起的在前）
func ListAppointments() []model.Appointment {
	store.mu.RLock()
	defer store.mu.RUnlock()
	result := make([]model.Appointment, len(store.appointments))
	copy(result, store.appointments)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// FindAppointment 按 ID 查询预约
func FindAppointment(id int64) (model.Appointment, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, appointment := range store.appointments {
		if appointment.ID == id {
			return appointment, true
		}
	}
	return model.Appointment{}, false
}

// CreateAppointment 追加一条新预约
func CreateAppointment(appointment model.Appointment) model.Appointment {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.nextID++
	appointment.ID = store.nextID
	store.appointments = append(store.appointments, appointment)
	return appointment
}

// UpdateAppointment 按 ID 覆盖更新预约，返回是否命中
func UpdateAppointment(appointment model.Appointment) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	for i := range store.appointments {
		if store.appointments[i].ID == appointment.ID {
			store.appointments[i] = appointment
			return true
		}
	}
	return false
}

// FindActiveConflict 查找某人在指定时段上的生效预约（pending/confirmed）。
// 命中说明该时段仍被占用，新申请必须被拒绝且原安排保留。
func FindActiveConflict(person, slot string) (model.Appointment, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, appointment := range store.appointments {
		if !constants.IsActiveAppointmentStatus(appointment.Status) || appointment.Slot != slot {
			continue
		}
		if appointment.Provider == person || appointment.Learner == person {
			return appointment, true
		}
	}
	return model.Appointment{}, false
}
