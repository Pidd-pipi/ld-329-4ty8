import { AppException } from '../errors/AppException';
import { logger } from '../logger/logger';
import type {
  Appointment,
  AppointmentActionPayload,
  CreateAppointmentPayload,
} from '../types/domain';

const API_BASE = '/api/appointments';

async function sendAppointmentRequest<T>(
  url: string,
  method: string,
  body?: unknown,
): Promise<T> {
  let response: Response;
  try {
    response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: body === undefined ? undefined : JSON.stringify(body),
    });
  } catch {
    throw new AppException('网络异常，预约操作未完成，请稍后重试');
  }

  if (response.ok) {
    return (await response.json()) as T;
  }

  let code = 'APPOINTMENT_ERROR';
  let message = '预约操作失败，请稍后重试';
  try {
    const data = (await response.json()) as { code?: string; message?: string };
    if (data.code) code = data.code;
    if (data.message) message = data.message;
  } catch {
    logger.warn('appointment error response without JSON body', response.status);
  }
  throw new AppException(message, code);
}

export function createAppointment(payload: CreateAppointmentPayload): Promise<Appointment> {
  return sendAppointmentRequest<Appointment>(API_BASE, 'POST', payload);
}

export function confirmAppointment(id: number, payload: AppointmentActionPayload): Promise<Appointment> {
  return sendAppointmentRequest<Appointment>(`${API_BASE}/${id}/confirm`, 'POST', payload);
}

export function completeAppointment(id: number, payload: AppointmentActionPayload): Promise<Appointment> {
  return sendAppointmentRequest<Appointment>(`${API_BASE}/${id}/complete`, 'POST', payload);
}

export function cancelAppointment(id: number, payload: AppointmentActionPayload): Promise<Appointment> {
  return sendAppointmentRequest<Appointment>(`${API_BASE}/${id}/cancel`, 'POST', payload);
}
