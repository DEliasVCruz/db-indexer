<script setup lang="ts">
import { computed } from "vue";
import { column } from "../globals/table";

const props = defineProps<{
  title: string;
  columnId: number;
}>();

const selected = computed(() => props.columnId === column.selected);
</script>

<template>
  <div
    class="mt-2 ml-0 w-1/4 border-b-2"
    @click.passive="column.select(props.columnId)"
  >
    <div class="w-[22rem] py-2 px-1 text-left after:h-3 after:bg-blue-200">
      {{ props.title }}
    </div>
    <span
      aria-hidden="true"
      class="block h-1 w-80 group-active:bg-blue-500"
      :class="{ 'bg-blue-500': selected }"
    ></span>
  </div>
</template>

<style scoped></style>
