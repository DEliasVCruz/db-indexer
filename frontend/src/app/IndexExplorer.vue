<script setup lang="ts">
import TheFooter from "@/components/TheFooter.vue";
import TheHeader from "@/components/TheHeader.vue";
import ToolBar from "@/components/ToolBar.vue";
import SideBar from "@/components/SideBar.vue";
import ResultTable from "@/components/ResultTable.vue";
import { mainContent } from "@/globals/content";

const contents = {
  ResultTable,
};
</script>

<template>
  <div class="outer-layer">
    <TheHeader :search-bar="true" />
    <div class="flex h-[85.5vh] w-screen flex-row">
      <SideBar />
      <main class="z-0 flex w-[94%] flex-initial flex-col pb-3 pt-1 pl-6">
        <ToolBar />
        <div class="h-[94%] overflow-x-auto">
          <component
            :is="contents[mainContent.current as keyof typeof contents]"
          ></component>
        </div>
      </main>
    </div>
    <TheFooter />
  </div>
</template>

<style scoped></style>
