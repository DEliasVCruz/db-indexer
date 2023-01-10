<script setup lang="ts">
import TableColumn from "@/components/TableColumn.vue";
import type { ColumnData } from "@/globals/types";

const props = defineProps<{
  columns: Array<ColumnData>;
}>();
</script>

<template>
  <div class="flex h-full min-w-max flex-row overflow-y-auto">
    <TableColumn
      v-for="column in props.columns"
      :key="column.name"
      :name="column.name"
      :cell-values="column.values"
    />
  </div>
</template>

<style scoped></style>
