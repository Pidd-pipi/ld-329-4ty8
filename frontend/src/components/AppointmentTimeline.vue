<template>
  <div class="appointment-list">
    <el-empty v-if="appointments.length === 0" description="还没有预约，去匹配卡片发起第一个吧" :image-size="72" />
    <el-timeline v-else>
      <el-timeline-item
        v-for="item in appointments"
        :key="item.id"
        :timestamp="timestampOf(item)"
        placement="top"
        :type="dotType(item.status)"
      >
        <article class="appointment-item" :class="{ 'is-done': item.status === 'cancelled' }">
          <div class="appointment-item__head">
            <strong>{{ item.pair }}</strong>
            <el-tag :type="APPOINTMENT_STATUS_TAG_TYPES[item.status]" size="small">
              {{ APPOINTMENT_STATUS_LABELS[item.status] }}
            </el-tag>
          </div>
          <p class="appointment-item__line">
            <el-tag size="small" effect="plain">{{ item.slot }}</el-tag>
            <el-tag size="small" type="info" effect="plain">{{ MEETING_MODE_LABELS[item.meetingMode] }}</el-tag>
          </p>
          <p class="muted">{{ item.offerSkill }} ↔ {{ item.wantedSkill }}</p>
          <p class="appointment-item__place">
            <template v-if="item.meetingMode === 'online'">会议链接：</template>
            <template v-else>线下地点：</template>
            <a v-if="item.meetingMode === 'online'" :href="item.location" target="_blank" rel="noopener">{{ item.location }}</a>
            <span v-else>{{ item.location }}</span>
          </p>
          <p v-if="item.agenda" class="muted">议程：{{ item.agenda }}</p>
          <ul class="appointment-item__meta">
            <li>发起人：{{ item.initiator }}</li>
            <li>
              确认人：
              <span v-if="item.confirmer">{{ item.confirmer }}</span>
              <span v-else class="muted">待对方确认</span>
            </li>
            <li v-if="item.status === 'cancelled' && item.cancelledBy">取消人：{{ item.cancelledBy }}</li>
          </ul>
          <div v-if="canAct(item)" class="appointment-item__actions">
            <el-button
              v-if="item.status === 'pending' && isOtherParty(item, currentUser)"
              type="primary"
              size="small"
              :loading="busyId === item.id"
              @click="emit('confirm', item)"
            >
              确认预约
            </el-button>
            <el-button
              v-if="item.status === 'confirmed'"
              type="success"
              size="small"
              :loading="busyId === item.id"
              @click="emit('complete', item)"
            >
              标记完成
            </el-button>
            <el-button
              v-if="item.status === 'pending' || item.status === 'confirmed'"
              type="danger"
              plain
              size="small"
              :loading="busyId === item.id"
              @click="emit('cancel', item)"
            >
              取消预约
            </el-button>
          </div>
        </article>
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script setup lang="ts">
import {
  APPOINTMENT_STATUS_LABELS,
  APPOINTMENT_STATUS_TAG_TYPES,
  MEETING_MODE_LABELS,
} from '../constants/appointment.constants';
import type { Appointment, AppointmentStatus } from '../types/domain';

defineProps<{
  appointments: Appointment[];
  currentUser: string;
  busyId?: number | null;
}>();

const emit = defineEmits<{
  (e: 'confirm', appointment: Appointment): void;
  (e: 'complete', appointment: Appointment): void;
  (e: 'cancel', appointment: Appointment): void;
}>();

function isOtherParty(item: Appointment, user: string): boolean {
  return (item.provider === user || item.learner === user) && item.initiator !== user;
}

function canAct(item: Appointment): boolean {
  return item.status === 'pending' || item.status === 'confirmed';
}

function dotType(status: AppointmentStatus): 'primary' | 'success' | 'info' | 'warning' {
  switch (status) {
    case 'pending':
      return 'warning';
    case 'confirmed':
      return 'primary';
    case 'completed':
      return 'success';
    default:
      return 'info';
  }
}

function timestampOf(item: Appointment): string {
  if (item.status === 'completed') return `完成于 ${item.completedAt ?? ''}`;
  if (item.status === 'cancelled') return `取消于 ${item.cancelledAt ?? ''}`;
  if (item.status === 'confirmed') return `确认于 ${item.confirmedAt ?? ''}`;
  return `发起于 ${item.createdAt}`;
}
</script>

<style scoped>
.appointment-item { padding: 12px 14px; border: 1px solid #e5e7eb; border-radius: 8px; background: #fff; }
.appointment-item.is-done { background: #fafafa; }
.appointment-item__head { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
.appointment-item__line { display: flex; flex-wrap: wrap; gap: 8px; margin: 8px 0; }
.appointment-item__place { margin: 6px 0; word-break: break-all; }
.appointment-item__place a { color: #2563eb; }
.appointment-item__meta { display: flex; flex-wrap: wrap; gap: 14px; margin: 8px 0 0; padding: 0; list-style: none; font-size: 13px; color: #475467; }
.appointment-item__actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 10px; }
</style>
