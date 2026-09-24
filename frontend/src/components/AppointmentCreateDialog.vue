<template>
  <el-dialog
    :model-value="modelValue"
    :title="`发起预约 · ${match?.provider ?? ''} × ${match?.learner ?? ''}`"
    width="460px"
    @update:model-value="(value: boolean) => !value && close()"
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="92px" @submit.prevent>
      <el-form-item label="交换内容">
        <span class="muted">{{ match?.offerSkill }} ↔ {{ match?.wantedSkill }}</span>
      </el-form-item>
      <el-form-item label="发起人">
        <el-tag size="small" effect="plain">{{ match?.learner }}</el-tag>
        <span class="muted confirm-hint">等待 {{ match?.provider }} 确认</span>
      </el-form-item>
      <el-form-item label="共同时段" prop="slot">
        <el-select v-model="form.slot" placeholder="选择双方共同时段" class="full-width">
          <el-option v-for="slot in match?.commonSlots ?? []" :key="slot" :label="slot" :value="slot" />
        </el-select>
      </el-form-item>
      <el-form-item label="约见方式" prop="locationType">
        <el-radio-group v-model="form.locationType">
          <el-radio-button value="ONLINE">线上会议</el-radio-button>
          <el-radio-button value="OFFLINE">线下地点</el-radio-button>
        </el-radio-group>
      </el-form-item>
      <el-form-item
        v-if="form.locationType === 'ONLINE'"
        label="会议链接"
        prop="place"
      >
        <el-input v-model="form.place" placeholder="https://meeting.example.com/room/..." clearable />
      </el-form-item>
      <el-form-item v-else label="线下地点" prop="place">
        <el-input v-model="form.place" placeholder="如：东校区图书馆三楼研讨间" clearable />
      </el-form-item>
      <el-form-item label="协商议程">
        <el-input
          v-model="form.agenda"
          type="textarea"
          :rows="2"
          maxlength="200"
          show-word-limit
          placeholder="想先练习什么、需要准备什么（选填）"
        />
      </el-form-item>
    </el-form>

    <el-alert
      v-if="conflictMessage"
      :title="conflictMessage"
      type="error"
      show-icon
      :closable="false"
      class="conflict-alert"
    />

    <template #footer>
      <el-button @click="close()">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="submit">发起预约</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import {
  APPOINTMENT_FORM_RULES,
  APPOINTMENT_TOAST_DURATION,
} from '../constants/appointment.constants';
import { logger } from '../logger/logger';
import { useAppointmentStore } from '../stores/appointment.store';
import type { CreateAppointmentPayload, LocationType, Match } from '../types/domain';
import { ApiError } from '../services/appointment.service';

const props = defineProps<{ modelValue: boolean; match: Match | null }>();
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>();

const store = useAppointmentStore();
const formRef = ref<FormInstance>();
const submitting = ref(false);
// 冲突提示会保留在对话框内（不关闭、不清空已填安排），普通错误仅 toast。
const conflictMessage = ref('');

const form = reactive<CreateAppointmentPayload>({
  matchId: 0,
  slot: '',
  locationType: 'ONLINE',
  place: '',
  agenda: '',
});

const onlineRule = (_: unknown, value: string, callback: (error?: Error) => void) => {
  if (!/^https?:\/\/\S+$/i.test(value.trim())) {
    callback(new Error(APPOINTMENT_FORM_RULES.meetingLinkInvalid));
    return;
  }
  callback();
};

const offlineRule = (_: unknown, value: string, callback: (error?: Error) => void) => {
  if (!value.trim()) {
    callback(new Error(APPOINTMENT_FORM_RULES.placeRequired));
    return;
  }
  callback();
};

const rules = computed<FormRules>(() => ({
  slot: [{ required: true, message: APPOINTMENT_FORM_RULES.slotRequired, trigger: 'change' }],
  locationType: [
    { required: true, message: APPOINTMENT_FORM_RULES.locationTypeRequired, trigger: 'change' },
  ],
  place: [
    {
      required: true,
      validator: form.locationType === 'ONLINE' ? onlineRule : offlineRule,
      trigger: 'blur',
    },
  ],
}));

watch(
  () => props.modelValue,
  (open) => {
    if (open && props.match) {
      form.matchId = props.match.id;
    }
  },
);

function resetForm() {
  form.slot = '';
  form.locationType = 'ONLINE' as LocationType;
  form.place = '';
  form.agenda = '';
  conflictMessage.value = '';
  formRef.value?.clearValidate();
}

function close() {
  emit('update:modelValue', false);
}

async function submit() {
  conflictMessage.value = '';
  if (!props.match) {
    return;
  }
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) {
    return;
  }
  submitting.value = true;
  try {
    const created = await store.create({ ...form });
    ElMessage.success({ message: `预约已发起，等待 ${created.confirmer} 确认`, duration: APPOINTMENT_TOAST_DURATION });
    close();
  } catch (error) {
    if (error instanceof ApiError && error.code === 'SLOT_CONFLICT') {
      // 明确提示冲突并保留原安排与当前已填写的内容，便于改选其他共同时段。
      conflictMessage.value = error.message;
    } else {
      const message = error instanceof Error ? error.message : '发起预约失败';
      ElMessage.error({ message, duration: APPOINTMENT_TOAST_DURATION });
      logger.warn('create appointment failed', error);
    }
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped>
.full-width { width: 100%; }
.confirm-hint { margin-left: 10px; font-size: 12px; }
.conflict-alert { margin-top: 4px; }
</style>
