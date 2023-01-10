<script setup lang="ts">
import SideBarOption from "./SideBarOption.vue";
import IconDirectory from "./IconDirectory.vue";
import IconIndex from "./IconIndex.vue";
import { reactive } from "vue";
import { mainSection } from "@/globals/section";

const props = defineProps<{
  show: boolean;
}>();

const options = reactive([
  { component: IconIndex, name: "Index" },
  { component: IconDirectory, name: "Dirs" },
]);

function navigateTo(section: string) {
  mainSection.setCurrent(section);
}
</script>

<template>
  <div
    class="left-0 z-10 m-0 w-16 flex-none bg-red-500 text-white shadow-lg"
    :class="{ hidden: !props.show }"
  >
    <div class="flex h-full flex-col p-2">
      <SideBarOption
        v-for="option in options"
        :key="option.name"
        :tooltip-text="option.name"
        @navigate="navigateTo"
      >
        <component :is="option.component"></component>
      </SideBarOption>
    </div>
  </div>
</template>

<style scoped></style>
