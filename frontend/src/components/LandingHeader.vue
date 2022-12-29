<script setup lang="ts"></script>

<template>
  <!-- <header class="absolute top-0 w-full bg-gray-900 pb-6"></header> -->
  <header
    class="relative h-72 w-screen overflow-hidden align-middle leading-none"
  >
    <svg
      data-name="Layer 1"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 1200 120"
      preserveAspectRatio="none"
      class="absolute inset-x-0 top-0 block h-72"
    >
      <path d="M1200 0L0 0 598.97 114.72 1200 0z" class="fill-black-500"></path>
    </svg>
    <svg
      data-name="Layer 1"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 1200 120"
      preserveAspectRatio="none"
      class="absolute inset-x-0 -top-6 block h-72"
    >
      <path d="M1200 0L0 0 598.97 114.72 1200 0z" class="fill-red-500"></path>
    </svg>
    <div class="inline-grid h-10 w-full grid-cols-3 grid-rows-3 md:h-40">
      <span></span>
      <span></span>
      <span></span>
      <span></span>
      <h1
        class="absoulte ms:text-4xl h-32 py-6 text-center text-3xl text-yellow-200"
      >
        Mailing Directory Indexer
      </h1>
    </div>
  </header>
</template>

<style scoped lang="postcss">
svg {
  width: calc(100% + 1.3px);
  z-index: -1;
}
</style>
