<script setup lang="ts">
import TableColumn from "@/components/TableColumn.vue";

export interface Props {
  columns?: { name: string; values: string[] }[];
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
  <div class="flex h-full w-full flex-row gap-0">
    <TableColumn
      v-for="column in props.columns"
      :key="column.name"
      :name="column.name"
      :cell-values="column.values"
    />
  </div>
</template>

<style scoped></style>
