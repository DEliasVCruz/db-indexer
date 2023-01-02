<script setup lang="ts">
import TheFooter from "@/components/TheFooter.vue";
import TheHeader from "@/components/TheHeader.vue";
import SideBar from "@/components/SideBar.vue";
import ResultTable from "@/components/ResultTable.vue";
</script>

<template>
  <div class="outer-layer">
    <TheHeader :search-bar="true" />
    <div class="flex h-[85.5vh] w-screen flex-row gap-4">
      <SideBar />
      <main class="w-[90vw]">
        <ResultTable />
      </main>
    </div>
    <TheFooter />
  </div>
</template>

<style scoped></style>
