<template>
  <FeatureCard :title="`${match.provider} × ${match.learner}`" :description="match.recommendation">
    <template #tag>
      <el-tag type="warning">{{ match.score }}%</el-tag>
    </template>
    <p class="muted">{{ match.offerSkill }} ↔ {{ match.wantedSkill }}</p>
    <div class="tag-row">
      <el-tag v-for="slot in match.commonSlots" :key="slot">{{ slot }}</el-tag>
    </div>
    <div class="card-actions">
      <el-button type="primary" size="small" @click="emit('create', match)">发起预约</el-button>
    </div>
  </FeatureCard>
</template>

<script setup lang="ts">
import FeatureCard from './FeatureCard.vue';
import type { Match } from '../types/domain';

defineProps<{ match: Match }>();
const emit = defineEmits<{
  (e: 'create', match: Match): void;
}>();
</script>

<style scoped>
.card-actions { display: flex; justify-content: flex-end; margin-top: 6px; }
</style>
