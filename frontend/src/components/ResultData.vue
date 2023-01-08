<script setup lang="ts">
import TableColumn from "@/components/TableColumn.vue";
import type { ColumnData } from "@/globals/types";

export interface Props {
  columns?: Array<ColumnData>;
}

const props = withDefaults(defineProps<Props>(), {
  columns: () => [
    { name: "From", values: ["Empty", "Empty", "Empty"] },
    { name: "Subject", values: ["Empty", "Empty", "Empty"] },
    { name: "Contents", values: ["Empty", "Empty", "Empty"] },
  ],
});
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
