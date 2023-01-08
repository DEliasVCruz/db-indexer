<script setup lang="ts">
import { reactive } from "vue";
import { results } from "@/globals/table";
import SearchField from "./SearchField.vue";

const props = defineProps<{
  searchText: string;
}>();

const searchFields = reactive({
  pagination: {
    from: results.from - 1,
    size: results.size,
  },
  queryData: {
    from: "",
    to: "",
    subject: "",
    contents: props.searchText,
  },
});
</script>

<template>
  <div
    class="absolute top-16 z-10 h-fit min-w-min max-w-max border-2 bg-white drop-shadow-lg"
  >
    <div
      class="flex w-[60vw] min-w-full max-w-[763px] flex-col gap-4 py-5 pl-5 pr-7"
    >
      <SearchField v-model="searchFields.queryData.from" :field="'From'" />
      <SearchField v-model="searchFields.queryData.to" :field="'To'" />
      <SearchField
        v-model="searchFields.queryData.subject"
        :field="'Subject'"
      />
      <div class="relative h-6">
        <button class="absolute right-0 h-8 w-20 p-1" @click.prevent>
          Search
        </button>
      </div>
    </div>
  </div>
</template>
