import type { AppointmentStatus, MeetingMode } from '../types/domain';

export const APPOINTMENT_STATUS_LABELS: Record<AppointmentStatus, string> = {
  pending: '待确认',
  confirmed: '待完成',
  completed: '已完成',
  cancelled: '已取消',
};

export const APPOINTMENT_STATUS_TAG_TYPES: Record<AppointmentStatus, 'warning' | 'primary' | 'success' | 'info'> = {
  pending: 'warning',
  confirmed: 'primary',
  completed: 'success',
  cancelled: 'info',
};

export const MEETING_MODE_LABELS: Record<MeetingMode, string> = {
  online: '线上会议',
  offline: '线下面对面',
};

// 生效中的预约状态：占用时段并计入首页预约数量
export const ACTIVE_APPOINTMENT_STATUSES: AppointmentStatus[] = ['pending', 'confirmed'];

// 与后端 errors/codes.go 保持一致
export const APPOINTMENT_ERROR_CODES = {
  slotConflict: 'SLOT_CONFLICT',
  invalidCommonSlot: 'INVALID_COMMON_SLOT',
  notConfirmer: 'NOT_CONFIRMER',
  invalidStatus: 'INVALID_STATUS',
} as const;

export const ACTOR_ALL = 'ALL';

export const APPOINTMENT_PANEL_TITLE = '预约确认';
export const ACTIVE_APPOINTMENT_METRIC_LABEL = '生效预约';
