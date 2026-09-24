<template>
  <el-dialog
    :model-value="modelValue"
    title="发起交换预约"
    width="460px"
    :close-on-click-modal="false"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <el-form
      v-if="match"
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="88px"
      @submit.prevent
    >
      <el-form-item label="交换匹配">
        <span class="muted">{{ match.provider }} × {{ match.learner }}（{{ match.offerSkill }} ↔ {{ match.wantedSkill }}）</span>
      </el-form-item>
      <el-form-item label="发起人" prop="initiator">
        <el-select v-model="form.initiator" placeholder="选择发起人">
          <el-option :label="match.provider" :value="match.provider" />
          <el-option :label="match.learner" :value="match.learner" />
        </el-select>
      </el-form-item>
      <el-form-item label="共同时段" prop="slot">
        <el-select v-model="form.slot" placeholder="选择双方共同时段">
          <el-option v-for="slot in match.commonSlots" :key="slot" :label="slot" :value="slot" />
        </el-select>
      </el-form-item>
      <el-form-item label="交换形式" prop="meetingMode">
        <el-radio-group v-model="form.meetingMode">
          <el-radio value="online">线上会议</el-radio>
          <el-radio value="offline">线下面对面</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item :label="form.meetingMode === 'online' ? '会议链接' : '线下地点'" prop="location">
        <el-input
          v-model="form.location"
          :placeholder="form.meetingMode === 'online' ? '请填写线上会议链接，如 https://meeting…' : '请填写线下地点，如东校区图书馆 302'"
          clearable
        />
      </el-form-item>
      <el-form-item label="交换议程">
        <el-input
          v-model="form.agenda"
          type="textarea"
          :rows="2"
          placeholder="可选：填写本次交换的具体安排"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:modelValue', false)">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">发起预约</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import { logger } from '../logger/logger';
import type { CreateAppointmentPayload, Match, MeetingMode } from '../types/domain';

const props = defineProps<{
  modelValue: boolean;
  match: Match | null;
  currentUser: string;
  submitting?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'submit', payload: CreateAppointmentPayload): void;
}>();

const formRef = ref<FormInstance>();

const createEmptyForm = () => ({
  initiator: '',
  slot: '',
  meetingMode: 'online' as MeetingMode,
  location: '',
  agenda: '',
});

const form = reactive(createEmptyForm());

const rules: FormRules = {
  initiator: [{ required: true, message: '请选择发起人', trigger: 'change' }],
  slot: [{ required: true, message: '请选择共同时段', trigger: 'change' }],
  meetingMode: [{ required: true, message: '请选择交换形式', trigger: 'change' }],
  location: [
    { required: true, message: '请填写会议链接或线下地点', trigger: 'blur' },
    {
      validator: (_rule, value: string, callback: (error?: Error) => void) => {
        if (form.meetingMode === 'online' && value && !/^https?:\/\//i.test(value.trim())) {
          callback(new Error('线上会议链接需以 http:// 或 https:// 开头'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    },
  ],
};

// 每次打开时按匹配与当前身份重置表单
watch(
  () => [props.modelValue, props.match] as const,
  ([visible, match]) => {
    if (!visible || !match) return;
    Object.assign(form, createEmptyForm());
    form.initiator = [match.provider, match.learner].includes(props.currentUser)
      ? props.currentUser
      : match.provider;
    form.slot = match.commonSlots[0] ?? '';
    formRef.value?.clearValidate();
  },
);

async function handleSubmit() {
  if (!props.match) return;
  try {
    await formRef.value?.validate();
  } catch {
    logger.warn('appointment form validation failed');
    return;
  }
  emit('submit', {
    matchId: props.match.id,
    initiator: form.initiator,
    slot: form.slot,
    meetingMode: form.meetingMode,
    location: form.location.trim(),
    agenda: form.agenda.trim(),
  });
}
</script>
