<script setup lang="ts">
import TheFooter from "../components/TheFooter.vue";
import TheHeader from "../components/TheHeader.vue";
</script>

<template>
  <div class="outer-layer">
    <TheHeader :side-bar="false" :search-bar="false" />
    <main>
      <h1>About</h1>
      <h2>The Usage</h2>
      <p></p>
      <h2>The Developer</h2>
      <p></p>
      <h2>The Motivation</h2>
      <p></p>
    </main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
