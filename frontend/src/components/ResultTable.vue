<script setup lang="ts">
import ResultData from "@/components/ResultData.vue";
import TableHeader from "@/components/TableHeader.vue";
import { columnData, row } from "@/globals/table";
import { onDeactivated } from "vue";

onDeactivated(() => {
  row.hover(0);
});
</script>

<template>
  <TableHeader :columns="columnData.columns" />
  <div class="h-[90%] w-full">
    <ResultData :columns="columnData.columns" />
  </div>
</template>
