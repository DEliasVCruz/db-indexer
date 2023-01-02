<script setup lang="ts">
import TableCell from "./TableCell.vue";
import TableColumnHeader from "./TableColumnHeader.vue";

export interface Props {
  name: string;
  cellValues?: string[];
}

const props = withDefaults(defineProps<Props>(), {
  cellValues: () => ["Empty"],
});
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
