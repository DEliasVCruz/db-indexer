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
    class="w-[22rem] truncate px-1 py-2 hover:bg-gray-100"
    :class="{
      'hovered-row': isHoveredRow,
      'text-center italic': !props.value,
      'normal-row': !isHoveredRow,
    }"
    @mouseover.passive="row.hover(props.rowId)"
    @mouseleave.passive="row.hover(0)"
    @dblclick.prevent="row.render(props.rowId - 1)"
  >
    {{ props.value ? props.value : "-" }}
  </p>
</template>
