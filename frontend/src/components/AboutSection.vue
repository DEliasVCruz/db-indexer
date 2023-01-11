<script setup lang="ts">
const props = defineProps<{
  title: string;
}>();
</script>

<template>
  <section class="mx-5 w-[88%] text-justify sm:ml-[22%] sm:w-[50%]">
    <h2 class="mb-1 text-xl font-semibold">{{ props.title }}</h2>
    <span aria-hidden="true" class="block h-0.5 w-full bg-gray-500"></span>
    <div class="mt-5 font-light">
      <slot />
    </div>
  </section>
</template>
