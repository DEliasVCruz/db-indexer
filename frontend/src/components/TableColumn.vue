<script setup lang="ts">
import TableCell from "./TableCell.vue";

const props = defineProps<{
  data: string[];
}>();
</script>

<template>
  <div class="flex grow cursor-pointer flex-col gap-0">
    <TableColumnHeader :title="props.name" />
    <div class="flex flex-col">
      <TableCell
        v-for="(cellValue, index) in props.cellValues"
        :key="index"
        :row-id="index + 1"
        :value="cellValue"
      />
    </div>
  </div>
</template>

<style scoped></style>
