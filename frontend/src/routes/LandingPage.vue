<script setup lang="ts">
import LandingHeader from "../components/LandingHeader.vue";
import TheFooter from "../components/TheFooter.vue";
import DropZone from "../components/DropZone.vue";
</script>

<template>
  <div class="outer-layer">
    <LandingHeader />
    <main>
      <DropZone />
    </main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
