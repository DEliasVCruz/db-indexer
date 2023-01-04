<script setup lang="ts">
import TableCell from "./TableCell.vue";

export interface Props {
  name: string;
  cellValues?: string[];
}

const props = withDefaults(defineProps<Props>(), {
  cellValues: () => ["Empty"],
});
</script>

<template>
  <div class="cursor-pointer">
    <TableCell
      v-for="(cellValue, index) in props.cellValues"
      :key="index"
      :row-id="index + 1"
      :value="cellValue"
    />
  </div>
</template>
