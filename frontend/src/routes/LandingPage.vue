<script setup lang="ts">
import LandingHeader from "../components/LandingHeader.vue";
import TheFooter from "../components/TheFooter.vue";
import DropZone from "../components/DropZone.vue";
</script>

<template>
  <div class="outer-layer">
    <LandingHeader />
    <main>
      <div class="flex flex-col items-center justify-center gap-2 pt-20">
        <DropZone class="mx-5 h-52 w-80 sm:w-96" />
        <span class="text-sm font-thin"
          >Allowed: gz, zip, tar | Max-Size: 500MB</span
        >
      </div>
    </main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
