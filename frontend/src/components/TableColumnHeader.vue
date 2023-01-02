<script setup lang="ts">
import { computed } from "vue";
import { column } from "../globals/directory";

const props = defineProps<{
  title: string;
}>();

const selected = computed(() => props.title === column.selected);
</script>

<template>
  <div
    class="mt-2 ml-0 flex flex-col border-b-2"
    @click.passive="column.select(props.title)"
  >
    <div class="py-2 px-1 text-left after:h-3 after:bg-blue-200">
      {{ props.title }}
    </div>
    <span
      aria-hidden="true"
      class="block h-1 w-full group-active:bg-blue-500"
      :class="{ 'bg-blue-500': selected }"
    ></span>
  </div>
</template>

<style scoped></style>
