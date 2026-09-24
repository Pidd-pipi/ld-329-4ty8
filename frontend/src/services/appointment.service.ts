import { AppException } from '../errors/AppException';
import { logger } from '../logger/logger';
import type {
  Appointment,
  AppointmentMutationResult,
  CreateAppointmentPayload,
} from '../types/domain';

const API_BASE = '/api';
const APPOINTMENTS_ENDPOINT = `${API_BASE}/appointments`;

// ApiError 携带后端返回的错误码（如 SLOT_CONFLICT），便于界面区分时段冲突。
export class ApiError extends Error {
  constructor(
    message: string,
    readonly code: string,
    readonly status: number,
  ) {
    super(message);
  }
}

async function requestAppointments<T>(url: string, init?: RequestInit): Promise<T> {
  let response: Response;
  try {
    response = await fetch(url, {
      headers: { 'Content-Type': 'application/json' },
      ...init,
    });
  } catch (networkError) {
    logger.error('appointment api unreachable', networkError);
    throw new AppException('网络异常，暂时无法连接预约服务');
  }

  if (!response.ok) {
    let code = 'APPOINTMENT_ERROR';
    let message = '预约操作失败，请稍后重试';
    try {
      const body = (await response.json()) as { code?: string; message?: string };
      code = body.code ?? code;
      message = body.message ?? message;
    } catch {
      logger.warn('appointment api returned non-json error', response.status);
    }
    throw new ApiError(message, code, response.status);
  }
  return response.json() as Promise<T>;
}

export function fetchAppointments(): Promise<Appointment[]> {
  return requestAppointments<Appointment[]>(APPOINTMENTS_ENDPOINT);
}

export function createAppointment(
  payload: CreateAppointmentPayload,
): Promise<AppointmentMutationResult> {
  return requestAppointments<AppointmentMutationResult>(APPOINTMENTS_ENDPOINT, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function confirmAppointment(id: number): Promise<AppointmentMutationResult> {
  return requestAppointments<AppointmentMutationResult>(`${APPOINTMENTS_ENDPOINT}/${id}/confirm`, {
    method: 'POST',
  });
}

export function completeAppointment(id: number): Promise<AppointmentMutationResult> {
  return requestAppointments<AppointmentMutationResult>(`${APPOINTMENTS_ENDPOINT}/${id}/complete`, {
    method: 'POST',
  });
}

export function cancelAppointment(id: number): Promise<AppointmentMutationResult> {
  return requestAppointments<AppointmentMutationResult>(`${APPOINTMENTS_ENDPOINT}/${id}/cancel`, {
    method: 'POST',
  });
}
