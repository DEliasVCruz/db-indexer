<script setup lang="ts">
import { reactive } from "vue";
import { columnData } from "@/globals/table";
import { searchText } from "@/lib/search";

function search() {
  searchText(searchInput.text).then((columns) => columnData.set(columns));
}

const searchInput = reactive({
  text: "",
  reset() {
    this.text = "";
  },
});
</script>

<template>
  <form
    class="flex h-full flex-row rounded-lg border-2 bg-white shadow-md"
    @submit.prevent
  >
    <button
      class="mr-1 rounded-full bg-gray-200 bg-opacity-0 p-3 hover:bg-opacity-60"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="h-6 w-6"
      >
        <path
          fill-rule="evenodd"
          d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
          clip-rule="evenodd"
        />
      </svg>
    </button>
    <input
      v-model="searchInput.text"
      class="my-2 w-[40vw] pr-2 outline-none"
      autocomplete="off"
      type="text"
      name="search-box"
      placeholder="Search index"
      @keyup.enter="search"
    />
    <button
      v-if="searchInput.text != ''"
      class="mr-1 rounded-full bg-gray-200 bg-opacity-0 p-3 hover:bg-opacity-60"
      @click="searchInput.reset()"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="h-6 w-6"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M6 18L18 6M6 6l12 12"
        />
      </svg>
    </button>
    <span v-else class="mr-3 w-10"></span>
    <button
      class="mr-1 rounded-full bg-gray-200 bg-opacity-0 p-3 hover:bg-opacity-60"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="h-5 w-5"
      >
        <path
          d="M10 3.75a2 2 0 10-4 0 2 2 0 004 0zM17.25 4.5a.75.75 0 000-1.5h-5.5a.75.75 0 000 1.5h5.5zM5 3.75a.75.75 0 01-.75.75h-1.5a.75.75 0 010-1.5h1.5a.75.75 0 01.75.75zM4.25 17a.75.75 0 000-1.5h-1.5a.75.75 0 000 1.5h1.5zM17.25 17a.75.75 0 000-1.5h-5.5a.75.75 0 000 1.5h5.5zM9 10a.75.75 0 01-.75.75h-5.5a.75.75 0 010-1.5h5.5A.75.75 0 019 10zM17.25 10.75a.75.75 0 000-1.5h-1.5a.75.75 0 000 1.5h1.5zM14 10a2 2 0 10-4 0 2 2 0 004 0zM10 16.25a2 2 0 10-4 0 2 2 0 004 0z"
        />
      </svg>
    </button>
  </form>
</template>
