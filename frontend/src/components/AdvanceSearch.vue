<script setup lang="ts">
import { reactive } from "vue";
import { columnData, results } from "@/globals/table";
import type { AdvanceSearch } from "@/globals/types";
import SearchField from "./SearchField.vue";
import { search } from "@/lib/search";
import { mainContent } from "@/globals/content";

const props = defineProps<{
  searchText: string;
}>();

const searchFields: AdvanceSearch = reactive({
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

function searchAdvance() {
  search(
    "advance",
    { advance: searchFields },
    "0",
    results.size.toString(),
    "contents"
  ).then((payload) => {
    results.setLastAdvanceQuery(searchFields);
    results.setLastQueryType("advance");
    results.setTotalResults(payload.total);
    results.resetRange();
    columnData.set(payload.columns);
    mainContent.setCurrent("ResultTable");
  });
}
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
        <button
          class="absolute right-0 h-8 w-20 p-1"
          @click.prevent="searchAdvance"
        >
          Search
        </button>
      </div>
    </div>
  </div>
</template>
