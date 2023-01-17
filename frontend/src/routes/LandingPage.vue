<script setup lang="ts">
import LandingHeader from "../components/LandingHeader.vue";
import TheFooter from "../components/TheFooter.vue";
import DropZone from "../components/DropZone.vue";
import { request } from "@/lib/http";

async function uploadFile(file: File) {
  const data = new FormData();
  data.append("file", file, file.name);
  const response = await request.put({
    endpoint: new URL("http://localhost:3000/api/index/indexName/upload"),
    dataTransfer: data,
  });

  const result = await response.json();
  alert(result.message);
}
</script>

<template>
  <div class="outer-layer">
    <LandingHeader />
    <main>
      <div class="flex flex-col items-center justify-center gap-2 pt-20">
        <DropZone class="mx-5 h-52 w-80 sm:w-96" @upload="uploadFile" />
        <span class="text-sm font-thin"
          >Allowed: tgz, zip, tar | Max-Size: 500MB</span
        >
      </div>
    </main>
    <TheFooter />
  </div>
</template>

<style scoped></style>
