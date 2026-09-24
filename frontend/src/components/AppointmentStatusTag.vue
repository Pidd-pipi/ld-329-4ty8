<template>
  <el-tag :type="APPOINTMENT_STATUS_TYPES[status]" effect="light" disable-transitions>
    {{ APPOINTMENT_STATUS_LABELS[status] }}
  </el-tag>
</template>

<script setup lang="ts">
import {
  APPOINTMENT_STATUS_LABELS,
  APPOINTMENT_STATUS_TYPES,
} from '../constants/appointment.constants';
import type { AppointmentStatus } from '../types/domain';

defineProps<{ status: AppointmentStatus }>();
</script>
