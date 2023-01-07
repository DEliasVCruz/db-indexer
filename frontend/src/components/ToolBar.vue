<script setup lang="ts">
import { columnData, results } from "@/globals/table";
import { searchText } from "@/lib/search";
import { computed } from "vue";

const isFirstPage = computed(() => {
  return results.page === 1;
});

const isLastPage = computed(() => {
  return results.to === results.total;
});

function nextPage() {
  if (isLastPage.value) {
    return;
  }
  results.nextPage();
  searchText(
    results.lastQuery,
    (results.from - 1).toString(),
    results.size.toString()
  ).then((payload) => {
    columnData.set(payload.columns);
  });
}

function prevPage() {
  if (isFirstPage.value) {
    return;
  }
  results.prevPage();
  searchText(
    results.lastQuery,
    (results.from - 1).toString(),
    results.size.toString()
  ).then((payload) => {
    columnData.set(payload.columns);
  });
}
</script>

<template>
  <div class="flex h-9 flex-row bg-green-300 py-1">
    <div class="mx-3 p-1.5 align-middle text-xs text-gray-700">
      {{ `${results.from} - ${results.to} of ${results.total}` }}
    </div>
    <button
      class="ml-0 rounded-full bg-gray-200 bg-opacity-0 pr-2 pl-1"
      :class="{
        'hover:bg-opacity-60': !isFirstPage,
        'cursor-default': isFirstPage,
      }"
      @click="prevPage"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="h-4 w-4"
        :class="{ 'stroke-gray-400': isFirstPage }"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M15.75 19.5L8.25 12l7.5-7.5"
        />
      </svg>
    </button>
    <button
      class="ml-0 rounded-full bg-gray-200 bg-opacity-0 pl-2 pr-1"
      :class="{
        'hover:bg-opacity-60': !isLastPage,
        'cursor-default': isLastPage,
      }"
      @click="nextPage"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="h-4 w-4"
        :class="{ 'stroke-gray-400': isLastPage }"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M8.25 4.5l7.5 7.5-7.5 7.5"
        />
      </svg>
    </button>
  </div>
</template>
