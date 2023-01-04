<script setup lang="ts">
import SideBarOption from "./SideBarOption.vue";
import IconDirectory from "./IconDirectory.vue";
import IconIndex from "./IconIndex.vue";
import { reactive } from "vue";

const options = reactive([
  { component: IconIndex, name: "Index" },
  { component: IconDirectory, name: "Dirs" },
]);
</script>

<template>
  <div
    class="left-0 z-10 m-0 flex h-full w-16 flex-none flex-col bg-red-500 text-white shadow-lg"
  >
    <SideBarOption
      v-for="option in options"
      :key="option.name"
      :tooltip-text="option.name"
    >
      <component :is="option.component"></component>
    </SideBarOption>
  </div>
</template>

<style scoped></style>
