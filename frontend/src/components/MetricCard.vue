<template>
  <article class="metric-card" :class="{ 'metric-card--highlight': highlight }">
    <span>{{ label }}</span>
    <strong>{{ value }}</strong>
  </article>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ label: string; value: number; highlight?: boolean }>(), {
  highlight: false,
});
</script>
