<script setup lang="ts">
import TableCell from "./TableCell.vue";

const props = defineProps<{
  data: string[];
}>();
</script>

<template>
  <div>
    <TableCell v-for="value in props.data" :key="props.data.indexOf(value)" />
  </div>
</template>

<style scoped></style>
