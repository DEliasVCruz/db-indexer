<script setup lang="ts">
import { row } from "@/globals/table";
</script>

<template>
  <div class="h-full w-full whitespace-pre-wrap">
    {{ row.data.get("Contents") }}
  </div>
</template>
