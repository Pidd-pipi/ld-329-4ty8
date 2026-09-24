<template>
  <div class="appointment-panel" v-loading="store.loading">
    <div class="appointment-panel__summary">
      <p class="muted">
        共 {{ store.appointments.length }} 条预约 ·
        <strong class="active-count">{{ store.activeCount }}</strong> 条进行中（待确认 / 待完成）
      </p>
    </div>

    <el-empty v-if="!store.loading && store.appointments.length === 0" description="还没有预约，去智能匹配卡片发起第一条吧" :image-size="70" />

    <AppointmentCard
      v-for="appointment in store.appointments"
      :key="appointment.id"
      :appointment="appointment"
      :acting="store.actingId === appointment.id"
      @confirm="onConfirm"
      @complete="onComplete"
      @cancel="onCancel"
    />
  </div>
</template>

<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus';
import { APPOINTMENT_TOAST_DURATION } from '../constants/appointment.constants';
import { logger } from '../logger/logger';
import { useAppointmentStore } from '../stores/appointment.store';
import AppointmentCard from './AppointmentCard.vue';
import { ApiError } from '../services/appointment.service';

const store = useAppointmentStore();

function handleError(error: unknown) {
  const message = error instanceof Error ? error.message : '操作失败，请稍后重试';
  ElMessage.error({ message, duration: APPOINTMENT_TOAST_DURATION });
  logger.warn('appointment action failed', error);
}

async function onConfirm(id: number) {
  try {
    await store.confirm(id);
    ElMessage.success({ message: '已确认，预约进入待完成', duration: APPOINTMENT_TOAST_DURATION });
  } catch (error) {
    handleError(error);
  }
}

async function onComplete(id: number) {
  try {
    await store.complete(id);
    ElMessage.success({ message: '交换已完成，感谢互助 🎉', duration: APPOINTMENT_TOAST_DURATION });
  } catch (error) {
    handleError(error);
  }
}

async function onCancel(id: number) {
  try {
    // 取消后该时段释放，双方均可再次发起预约。
    await ElMessageBox.confirm(
      '取消后该共同时段将重新开放预约，确定取消这次约见吗？',
      '取消预约',
      { confirmButtonText: '确定取消', cancelButtonText: '再想想', type: 'warning' },
    );
  } catch {
    return; // 用户放弃取消
  }
  try {
    await store.cancel(id);
    ElMessage.info({ message: '预约已取消，时段可再次预约', duration: APPOINTMENT_TOAST_DURATION });
  } catch (error) {
    if (error instanceof ApiError) {
      handleError(error);
    } else {
      logger.warn('cancel aborted or failed', error);
    }
  }
}
</script>

<style scoped>
.appointment-panel__summary { margin-bottom: 10px; }
.active-count { color: #409eff; font-size: 15px; }
</style>
