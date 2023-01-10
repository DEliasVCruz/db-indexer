<script setup lang="ts">
const props = defineProps<{
  tooltipText: string;
}>();

const emit = defineEmits<{
  (e: "navigate", section: string): void;
}>();
</script>

<template>
  <div
    class="sidebar-option group"
    @click.prevent="emit('navigate', props.tooltipText)"
  >
    <slot />
    <span class="sidebar-tooltip group-hover:scale-100">
      {{ props.tooltipText }}
    </span>
  </div>
</template>

<style scoped></style>
