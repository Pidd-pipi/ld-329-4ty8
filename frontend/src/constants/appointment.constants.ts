import type { AppointmentStatus, LocationType } from '../types/domain';

// 预约状态的中文标签与标签颜色，与后端状态机保持一致。
export const APPOINTMENT_STATUS_LABELS: Record<AppointmentStatus, string> = {
  PENDING: '待对方确认',
  CONFIRMED: '待完成',
  COMPLETED: '已完成',
  CANCELLED: '已取消',
};

export const APPOINTMENT_STATUS_TYPES: Record<AppointmentStatus, 'warning' | 'primary' | 'success' | 'info'> = {
  PENDING: 'warning',
  CONFIRMED: 'primary',
  COMPLETED: 'success',
  CANCELLED: 'info',
};

// 仍占用时段的状态：待确认与待完成，已完成或已取消后时段重新开放。
export const ACTIVE_APPOINTMENT_STATUSES: AppointmentStatus[] = ['PENDING', 'CONFIRMED'];

export const LOCATION_TYPE_LABELS: Record<LocationType, string> = {
  ONLINE: '线上会议',
  OFFLINE: '线下地点',
};

export const isActiveAppointment = (status: AppointmentStatus): boolean =>
  ACTIVE_APPOINTMENT_STATUSES.includes(status);

// 预约表单校验文案。
export const APPOINTMENT_FORM_RULES = {
  slotRequired: '请选择一个双方共同时段',
  locationTypeRequired: '请选择线上会议或线下地点',
  meetingLinkInvalid: '请填写有效的线上会议链接（需以 http:// 或 https:// 开头）',
  placeRequired: '请填写线下约见地点',
} as const;

// Element Plus 消息提示时长（毫秒）。
export const APPOINTMENT_TOAST_DURATION = 3200;
