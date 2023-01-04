<script setup lang="ts">
import TheFooter from "@/components/TheFooter.vue";
import TheHeader from "@/components/TheHeader.vue";
import SideBar from "@/components/SideBar.vue";
import ResultTable from "@/components/ResultTable.vue";
import TableHeader from "@/components/TableHeader.vue";
import { columnData } from "@/globals/table";
</script>

<template>
  <div class="outer-layer">
    <TheHeader :search-bar="true" />
    <div class="flex h-[85.5vh] w-screen flex-row">
      <SideBar />
      <main class="z-0 w-[94%] flex-initial overflow-x-auto pb-3 pt-1 pl-6">
        <!-- <ResultTable /> -->
        <TableHeader :columns="columnData.columns" />
        <ResultTable :columns="columnData.columns" />
      </main>
    </div>
    <TheFooter />
  </div>
</template>

<style scoped></style>
