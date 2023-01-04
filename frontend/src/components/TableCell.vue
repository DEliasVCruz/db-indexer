<script setup lang="ts">
import { computed } from "vue";
import { row } from "../globals/table";

const props = defineProps<{
  value: string;
  rowId: number;
}>();

const isHoveredRow = computed(() => props.rowId === row.hovered);
</script>

<template>
  <p
    class="truncate hover:bg-gray-100"
    :class="{ 'hovered-row': isHoveredRow, 'normal-row': !isHoveredRow }"
    @mouseover.passive="row.hover(props.rowId)"
    @mouseleave.passive="row.hover(0)"
  >
    {{ props.value }}
  </p>
</template>
