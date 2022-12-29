<script setup lang="ts">
import TheSearchBar from "../components/TheSearchBar.vue";

const props = defineProps<{
  searchBar: boolean;
}>();
</script>

<template>
  <div class="relative">
    <header class="fixed top-0 left-0 right-0 h-14 bg-blue-900 py-6">
      <div class="absolute top-0 left-0 flex h-14 flex-row gap-x-5">
        <div class="flex flex-row">
          <img
            class="h-14 w-14"
            src="../assets/menu.svg"
            alt="icon to toggle nav bar"
          />
          <span class="-ml-1 mb-1 inline-flex items-center gap-x-1">
            <img
              class="hidden h-6 w-6 fill-none sm:inline-block"
              src="../assets/logo.svg"
              alt="page logo in the shape of an @ symbol"
            />
            <h1>MailInDexer</h1>
          </span>
        </div>
        <div v-if="props.searchBar">
          <TheSearchBar />
        </div>
      </div>
    </header>
  </div>
</template>

<style scoped lang="postcss">
svg {
  width: calc(100% + 1.3px);
  z-index: -1;
}
</style>
