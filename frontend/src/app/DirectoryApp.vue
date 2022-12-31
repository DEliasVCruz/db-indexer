<script setup lang="ts">
import TheFooter from "@/components/TheFooter.vue";
import TheHeader from "@/components/TheHeader.vue";
import SideBar from "@/components/SideBar.vue";
import ResultTable from "@/components/ResultTable.vue";
</script>

<template>
  <div class="outer-layer">
    <TheHeader :search-bar="true" />
    <main flex flex-row>
      <SideBar />
      <ResultTable />
    </main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
