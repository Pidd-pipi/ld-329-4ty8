<template>
  <main class="page-shell" v-loading="loading">
    <AppHeader :unread="overview?.metrics.unread ?? 0" />

    <section v-if="overview" class="metrics-grid">
      <MetricCard label="已发布技能" :value="overview.metrics.skills" />
      <MetricCard label="活跃需求" :value="overview.metrics.needs" />
      <MetricCard label="智能匹配" :value="overview.metrics.matches" />
      <MetricCard label="生效预约" :value="overview.metrics.appointments" />
      <MetricCard label="评价记录" :value="overview.metrics.reviews" />
    </section>

    <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" style="margin-bottom: 14px" />

    <section v-if="overview" class="workspace-grid">
      <div class="panel">
        <h2>技能发布</h2>
        <FeatureCard v-for="skill in overview.skills" :key="skill.id" :title="skill.title" :description="skill.description">
          <template #tag><el-tag>{{ skill.category }} {{ skill.level }}%</el-tag></template>
          <div class="tag-row">
            <el-tag v-for="slot in skill.timeSlots" :key="slot" effect="plain">{{ slot }}</el-tag>
            <el-tag v-for="reward in skill.rewards" :key="reward" type="success" effect="plain">{{ reward }}</el-tag>
          </div>
          <small>{{ skill.owner }} · {{ skill.campus }} · {{ skill.portfolio }}</small>
        </FeatureCard>
      </div>

      <div class="panel">
        <h2>需求浏览</h2>
        <el-table :data="overview.needs" size="small">
          <el-table-column prop="title" label="需求" min-width="170" />
          <el-table-column prop="category" label="类别" width="82" />
          <el-table-column prop="campus" label="校区" width="96" />
          <el-table-column prop="responses" label="响应" width="72" sortable />
        </el-table>
      </div>

      <div class="panel">
        <h2>智能匹配</h2>
        <MatchCard v-for="match in overview.matches" :key="match.id" :match="match" @create="openCreateDialog" />
      </div>

      <div class="panel">
        <div class="panel__head">
          <h2>预约确认</h2>
          <ActorSwitch :current-user="currentUser" :users="actorUsers" @change="currentUser = $event" />
        </div>
        <AppointmentTimeline
          :appointments="overview.appointments"
          :current-user="currentUser"
          :busy-id="busyAppointmentId"
          @confirm="handleConfirm"
          @complete="handleComplete"
          @cancel="handleCancel"
        />
      </div>

      <div class="panel profile-panel">
        <div>
          <h2>个人主页与技能墙</h2>
          <h3>{{ overview.profile.name }}</h3>
          <p>{{ overview.profile.major }} · {{ overview.profile.creditLevel }}</p>
          <el-progress :percentage="overview.profile.creditScore" />
          <ul>
            <li v-for="item in overview.profile.history" :key="item">{{ item }}</li>
          </ul>
        </div>
        <RadarChart :radar="overview.profile.radar" />
      </div>

      <div class="panel">
        <h2>评价信用</h2>
        <FeatureCard v-for="review in overview.reviews" :key="review.id" :title="`${review.from} → ${review.to}`" :description="review.content">
          <template #tag><el-rate :model-value="review.rating" disabled size="small" /></template>
        </FeatureCard>
      </div>

      <div class="panel">
        <h2>消息通知</h2>
        <FeatureCard v-for="conversation in overview.messages" :key="conversation.id" :title="conversation.withUser" :description="conversation.messages.join(' / ')">
          <template #tag><el-badge :value="conversation.unread" /></template>
        </FeatureCard>
      </div>
    </section>

    <AppointmentDialog
      v-model="dialogVisible"
      :match="activeMatch"
      :current-user="currentUser"
      :submitting="creating"
      @submit="handleCreate"
    />
  </main>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import AppHeader from '../components/AppHeader.vue';
import FeatureCard from '../components/FeatureCard.vue';
import MetricCard from '../components/MetricCard.vue';
import RadarChart from '../components/RadarChart.vue';
import MatchCard from '../components/MatchCard.vue';
import AppointmentDialog from '../components/AppointmentDialog.vue';
import AppointmentTimeline from '../components/AppointmentTimeline.vue';
import ActorSwitch from '../components/ActorSwitch.vue';
import { AppException } from '../errors/AppException';
import { logger } from '../logger/logger';
import { fetchOverview } from '../services/storage.service';
import {
  cancelAppointment,
  completeAppointment,
  confirmAppointment,
  createAppointment,
} from '../services/appointment.service';
import type {
  Appointment,
  CreateAppointmentPayload,
  Match,
  Overview,
} from '../types/domain';

const overview = ref<Overview | null>(null);
const loading = ref(true);
const error = ref('');

// 演示环境中切换操作身份，用于验证发起人/对方各自的确认、取消权限
const currentUser = ref('林澈');
const actorUsers = computed(() => {
  if (!overview.value) return [currentUser.value];
  const users = new Set<string>([currentUser.value]);
  overview.value.matches.forEach((match) => {
    users.add(match.provider);
    users.add(match.learner);
  });
  return Array.from(users);
});

const dialogVisible = ref(false);
const creating = ref(false);
const activeMatch = ref<Match | null>(null);
const busyAppointmentId = ref<number | null>(null);

async function loadOverview() {
  overview.value = await fetchOverview();
}

function openCreateDialog(match: Match) {
  activeMatch.value = match;
  dialogVisible.value = true;
}

function reportFailure(err: unknown, fallback: string) {
  const message = err instanceof AppException ? err.message : fallback;
  error.value = '';
  logger.error('appointment operation failed', err);
  ElMessage.error(message);
}

async function runAppointmentAction(
  appointmentId: number,
  action: () => Promise<Appointment>,
  successMessage: string,
) {
  busyAppointmentId.value = appointmentId;
  try {
    const result = await action();
    await loadOverview();
    ElMessage.success(successMessage);
    return result;
  } catch (err) {
    reportFailure(err, '预约操作失败，请稍后重试');
    throw err;
  } finally {
    busyAppointmentId.value = null;
  }
}

async function handleCreate(payload: CreateAppointmentPayload) {
  creating.value = true;
  try {
    await createAppointment(payload);
    await loadOverview();
    dialogVisible.value = false;
    ElMessage.success('预约已发起，等待对方确认');
  } catch (err) {
    // 时段冲突等错误已由后端返回明确文案；保留弹窗与原安排，允许调整时段后重试
    reportFailure(err, '发起预约失败，请稍后重试');
  } finally {
    creating.value = false;
  }
}

function handleConfirm(appointment: Appointment) {
  void runAppointmentAction(
    appointment.id,
    () => confirmAppointment(appointment.id, { actor: currentUser.value }),
    '已确认，预约进入待完成',
  ).catch(() => undefined);
}

function handleComplete(appointment: Appointment) {
  void runAppointmentAction(
    appointment.id,
    () => completeAppointment(appointment.id, { actor: currentUser.value }),
    '交换已完成，感谢你的互助分享',
  ).catch(() => undefined);
}

async function handleCancel(appointment: Appointment) {
  try {
    await ElMessageBox.confirm(
      `取消后「${appointment.slot}」时段将释放，双方都可再次预约。确认取消吗？`,
      '取消预约',
      { type: 'warning', confirmButtonText: '确认取消', cancelButtonText: '再想想' },
    );
  } catch {
    return;
  }
  try {
    await runAppointmentAction(
      appointment.id,
      () => cancelAppointment(appointment.id, { actor: currentUser.value }),
      '预约已取消，时段可再次预约',
    );
  } catch {
    /* 错误已统一提示 */
  }
}

onMounted(async () => {
  try {
    await loadOverview();
  } catch (err) {
    error.value = err instanceof Error ? err.message : '加载失败';
  } finally {
    loading.value = false;
  }
});
</script>
