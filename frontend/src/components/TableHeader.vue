<script setup lang="ts">
import type { ColumnData } from "@/globals/types";
import TableColumnHeader from "./TableColumnHeader.vue";

const props = defineProps<{
  columns: Array<ColumnData>;
}>();
</script>

<template>
  <div class="flex min-w-max cursor-pointer flex-row">
    <TableColumnHeader
      v-for="column in props.columns"
      :key="column.name"
      :title="column.name"
    />
  </div>
</template>
