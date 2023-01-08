<script setup lang="ts">
import type { ColumnData } from "@/globals/types";
import TableColumnHeader from "./TableColumnHeader.vue";

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
  <div class="flex min-w-max cursor-pointer flex-row">
    <TableColumnHeader
      v-for="column in props.columns"
      :key="column.name"
      :title="column.name"
    />
  </div>
</template>
