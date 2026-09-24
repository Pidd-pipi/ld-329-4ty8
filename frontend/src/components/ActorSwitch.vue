<template>
  <div class="actor-switch">
    <span class="muted">当前操作身份</span>
    <el-select :model-value="currentUser" size="small" style="width: 132px" @update:model-value="emit('change', String($event))">
      <el-option v-for="user in users" :key="user" :label="user" :value="user" />
    </el-select>
  </div>
</template>

<script setup lang="ts">
defineProps<{ currentUser: string; users: string[] }>();
const emit = defineEmits<{ (e: 'change', user: string): void }>();
</script>

<style scoped>
.actor-switch { display: flex; align-items: center; gap: 8px; font-size: 13px; }
</style>
