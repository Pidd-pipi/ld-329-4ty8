<template>
  <article class="appointment-card" :class="`appointment-card--${statusClass}`">
    <header class="appointment-card__header">
      <div class="appointment-card__people">
        <strong>{{ appointment.pair }}</strong>
        <AppointmentStatusTag :status="appointment.status" />
      </div>
      <small>{{ appointment.offerSkill }} ↔ {{ appointment.wantedSkill }}</small>
    </header>

    <ul class="appointment-card__meta">
      <li><span class="meta-label">共同时段</span>{{ appointment.slot }}</li>
      <li>
        <span class="meta-label">{{ locationLabel }}</span>
        <el-link
          v-if="appointment.locationType === 'ONLINE'"
          :href="appointment.place"
          target="_blank"
          rel="noopener"
          type="primary"
          class="meeting-link"
        >{{ appointment.place }}</el-link>
        <template v-else>{{ appointment.place }}</template>
      </li>
      <li v-if="appointment.agenda"><span class="meta-label">协商议程</span>{{ appointment.agenda }}</li>
      <li><span class="meta-label">发起人</span>{{ appointment.initiator }}</li>
      <li>
        <span class="meta-label">确认人</span>
        <template v-if="appointment.status === 'PENDING'">待 {{ appointment.confirmer }} 确认</template>
        <template v-else>{{ appointment.confirmer }}</template>
      </li>
      <li><span class="meta-label">更新时间</span>{{ appointment.updatedAt }}</li>
    </ul>

    <footer v-if="isActive" class="appointment-card__actions">
      <el-button
        v-if="appointment.status === 'PENDING'"
        type="primary"
        size="small"
        :loading="acting"
        @click="emit('confirm', appointment.id)"
      >{{ appointment.confirmer }} 确认预约</el-button>
      <el-button
        v-if="appointment.status === 'CONFIRMED'"
        type="success"
        size="small"
        :loading="acting"
        @click="emit('complete', appointment.id)"
      >标记完成</el-button>
      <el-button size="small" :loading="acting" @click="emit('cancel', appointment.id)">
        取消预约
      </el-button>
    </footer>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import {
  LOCATION_TYPE_LABELS,
  isActiveAppointment,
} from '../constants/appointment.constants';
import type { Appointment } from '../types/domain';
import AppointmentStatusTag from './AppointmentStatusTag.vue';

const props = defineProps<{
  appointment: Appointment;
  acting: boolean;
}>();

const emit = defineEmits<{
  confirm: [id: number];
  complete: [id: number];
  cancel: [id: number];
}>();

const isActive = computed(() => isActiveAppointment(props.appointment.status));
const locationLabel = computed(() => LOCATION_TYPE_LABELS[props.appointment.locationType]);
const statusClass = computed(() => props.appointment.status.toLowerCase());
</script>

<style scoped>
.appointment-card { background: #fff; border: 1px solid #e5e7eb; border-left: 4px solid #d0d5dd; border-radius: 8px; padding: 14px 16px; margin-bottom: 12px; }
.appointment-card--pending { border-left-color: #e6a23c; }
.appointment-card--confirmed { border-left-color: #409eff; }
.appointment-card--completed { border-left-color: #67c23a; }
.appointment-card--cancelled { opacity: 0.72; }
.appointment-card__header { display: flex; flex-direction: column; gap: 4px; margin-bottom: 8px; }
.appointment-card__people { display: flex; align-items: center; gap: 10px; }
.appointment-card__meta { list-style: none; margin: 0; padding: 0; display: grid; gap: 6px; font-size: 13px; }
.meta-label { display: inline-block; width: 64px; color: #667085; }
.meeting-link { max-width: 100%; overflow-wrap: anywhere; }
.appointment-card__actions { margin-top: 12px; display: flex; gap: 8px; flex-wrap: wrap; }
</style>
