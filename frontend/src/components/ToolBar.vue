<script setup lang="ts">
import { columnData, results } from "@/globals/table";
import { search } from "@/lib/search";
import { computed } from "vue";
import { mainContent } from "@/globals/content";

const isFirstPage = computed(() => {
  return results.page === 1;
});

const isLastPage = computed(() => {
  return results.to === results.total;
});

const isMailView = computed(() => {
  return mainContent.current === "MailView";
});

const isResultTable = computed(() => {
  return mainContent.current === "ResultTable";
});

function nextPage() {
  if (isLastPage.value) {
    return;
  }
  results.nextPage();
  search(
    results.lastQueryType,
    results.lastQuery,
    results.from,
    results.size,
    "contents"
  ).then((payload) => {
    columnData.set(payload.columns);
  });
}

function prevPage() {
  if (isFirstPage.value) {
    return;
  }
  results.prevPage();
  search(
    results.lastQueryType,
    results.lastQuery,
    results.from,
    results.size,
    "contents"
  ).then((payload) => {
    columnData.set(payload.columns);
  });
}
</script>

<template>
  <div class="flex h-9 flex-row py-1">
    <button
      class="ml-2 rounded-full bg-gray-200 bg-opacity-0 pr-1 pl-1 hover:bg-opacity-60"
      :class="{ hidden: !isMailView }"
      @click.passive="mainContent.setCurrent('ResultTable')"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="h-5 w-5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M19.5 12h-15m0 0l6.75 6.75M4.5 12l6.75-6.75"
        />
      </svg>
    </button>
    <div
      class="ml-auto mr-5 flex w-fit flex-row"
      :class="{ hidden: !isResultTable }"
    >
      <div class="mx-3 p-1.5 align-middle text-xs text-gray-700">
        {{ `${results.from} - ${results.to} of ${results.total}` }}
      </div>
      <button
        class="rounded-full bg-gray-200 bg-opacity-0 pr-2 pl-1"
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
        class="rounded-full bg-gray-200 bg-opacity-0 pl-2 pr-1"
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
  </div>
</template>
