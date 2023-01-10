<script setup lang="ts">
import { row } from "@/globals/table";
</script>

<template>
  <div class="flex flex-row">
    <div class="w-[3%]"></div>
    <p class="h-full w-full whitespace-pre-wrap">
      {{ row.data.get("Contents") }}
    </p>
  </div>
</template>
