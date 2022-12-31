<script setup lang="ts">
import TheFooter from "../components/TheFooter.vue";
import TheHeader from "../components/TheHeader.vue";
</script>

<template>
  <div class="outer-layer">
    <TheHeader :search-bar="false" />
    <main></main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
