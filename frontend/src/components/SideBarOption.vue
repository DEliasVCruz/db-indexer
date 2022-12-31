<script setup lang="ts">
const props = defineProps<{
  tooltipText: string;
}>();
</script>

<template>
  <div class="sidebar-option group">
    <slot />
    <span class="sidebar-tooltip group-hover:scale-100">
      {{ props.tooltipText }}
    </span>
  </div>
</template>

<style scoped></style>
