<script setup lang="ts">
import TheFooter from "@/components/TheFooter.vue";
import TheHeader from "@/components/TheHeader.vue";
import ToolBar from "@/components/ToolBar.vue";
import SideBar from "@/components/SideBar.vue";
import ResultTable from "@/components/ResultTable.vue";
import MailView from "@/components/MailView.vue";
import NoContent from "@/components/NoContent.vue";
import { mainContent } from "@/globals/content";
import { ref } from "vue";

const showSidebar = ref(true);

const contents = {
  ResultTable,
  MailView,
  NoContent,
};
</script>

<template>
  <div class="outer-layer">
    <TheHeader
      :side-bar="true"
      :search-bar="true"
      @toggle="showSidebar = !showSidebar"
    />
    <div class="flex h-[85.5vh] w-screen flex-row gap-1">
      <SideBar :show="showSidebar" />
      <main
        class="z-0 m-4 flex flex-initial flex-col rounded-3xl border-2 p-2"
        :class="{ 'w-[98%]': !showSidebar, 'w-[94%]': showSidebar }"
      >
        <ToolBar />
        <div class="h-[94%] overflow-x-auto">
          <KeepAlive :max="2">
            <component
              :is="contents[mainContent.current as keyof typeof contents]"
            ></component>
          </KeepAlive>
        </div>
      </main>
    </div>
    <TheFooter />
  </div>
</template>

<style scoped></style>
