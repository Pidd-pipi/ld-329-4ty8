package repository

import (
	"sync"

	"cyskillswap/internal/model"
)

// appointmentStore 预约数据的内存存储，演示项目中替代数据库表。
// 预约状态会被发起、确认、完成、取消等写操作改变，因此需要互斥保护。
type appointmentStore struct {
	mu           sync.RWMutex
	sequence     int
	appointments []model.Appointment
}

var appointmentData = newAppointmentStore()

func newAppointmentStore() *appointmentStore {
	store := &appointmentStore{}
	store.seed()
	return store
}

// read 加读锁执行只读操作。
func (s *appointmentStore) read(fn func() any) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fn()
}

// write 加写锁执行变更操作。
func (s *appointmentStore) write(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn()
}

// nextID 在写锁内分配自增主键。
func (s *appointmentStore) nextID() int {
	s.sequence++
	return s.sequence
}

// add 在写锁内追加预约。
func (s *appointmentStore) add(appointment model.Appointment) {
	s.appointments = append(s.appointments, appointment)
}

// replace 在写锁内按编号替换预约。
func (s *appointmentStore) replace(appointment model.Appointment) bool {
	for i := range s.appointments {
		if s.appointments[i].ID == appointment.ID {
			s.appointments[i] = appointment
			return true
		}
	}
	return false
}

// snapshot 返回预约列表的浅拷贝，避免调用方绕过锁修改底层切片。
func (s *appointmentStore) snapshot() []model.Appointment {
	items := make([]model.Appointment, len(s.appointments))
	copy(items, s.appointments)
	return items
}

// findByID 按编号查找预约。
func (s *appointmentStore) findByID(id int) (model.Appointment, bool) {
	for _, appointment := range s.appointments {
		if appointment.ID == id {
			return appointment, true
		}
	}
	return model.Appointment{}, false
}
