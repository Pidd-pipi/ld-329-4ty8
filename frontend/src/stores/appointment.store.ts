import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import { isActiveAppointment } from '../constants/appointment.constants';
import {
  cancelAppointment,
  completeAppointment,
  confirmAppointment,
  createAppointment,
  fetchAppointments,
} from '../services/appointment.service';
import type { Appointment, CreateAppointmentPayload } from '../types/domain';

// useAppointmentStore 持有预约区全部状态，发起、确认、完成、取消后立即替换列表，
// 首页预约数量通过 activeCount 计算属性同步更新。
export const useAppointmentStore = defineStore('appointment', () => {
  const appointments = ref<Appointment[]>([]);
  const loaded = ref(false);
  const loading = ref(false);
  const actingId = ref<number | null>(null);

  const activeCount = computed(() =>
    appointments.value.filter((item) => isActiveAppointment(item.status)).length,
  );

  function replaceAppointments(items: Appointment[]) {
    appointments.value = items;
    loaded.value = true;
  }

  async function load(force = false) {
    if (loaded.value && !force) {
      return;
    }
    loading.value = true;
    try {
      replaceAppointments(await fetchAppointments());
    } finally {
      loading.value = false;
    }
  }

  // seedFromOverview 使用首页概览数据做首次渲染，随后再由预约区接口校准。
  function seedFromOverview(items: Appointment[]) {
    if (!loaded.value) {
      replaceAppointments(items);
    }
  }

  async function create(payload: CreateAppointmentPayload) {
    const result = await createAppointment(payload);
    replaceAppointments(result.appointments);
    return result.appointment;
  }

  async function confirm(id: number) {
    actingId.value = id;
    try {
      const result = await confirmAppointment(id);
      replaceAppointments(result.appointments);
    } finally {
      actingId.value = null;
    }
  }

  async function complete(id: number) {
    actingId.value = id;
    try {
      const result = await completeAppointment(id);
      replaceAppointments(result.appointments);
    } finally {
      actingId.value = null;
    }
  }

  async function cancel(id: number) {
    actingId.value = id;
    try {
      const result = await cancelAppointment(id);
      replaceAppointments(result.appointments);
    } finally {
      actingId.value = null;
    }
  }

  return {
    appointments,
    loaded,
    loading,
    actingId,
    activeCount,
    load,
    seedFromOverview,
    create,
    confirm,
    complete,
    cancel,
  };
});
