<script setup lang="ts">
import { reactive } from "vue";

function handleInput(e: Event) {
  const files = (e.target as HTMLInputElement).files as FileList;
  emit("upload", [...files]);
}

function handleDrop(e: Event) {
  hover.toggle();
  const files = (e as DragEvent).dataTransfer?.files as FileList;
  emit("upload", [...files]);
}

const hover = reactive({
  is: false,
  toggle() {
    this.is = !this.is;
  },
});

const emit = defineEmits<{
  (e: "upload", files: File[]): void;
}>();
</script>

<template>
  <div
    class="dropzone"
    :class="
      hover.is ? 'bg-gray-100 bg-opacity-60' : 'bg-gray-100 bg-opacity-10'
    "
    v-bind="$attrs"
    @drop.prevent="handleDrop"
    @dragenter.prevent="hover.toggle"
    @dragleave.prevent="hover.toggle"
    @dragover.prevent
  >
    <span class="font-light">Drag and Drop</span>
    <span class="font-light">OR</span>
    <label
      for="fileInput"
      class="input-button"
      @mouseover="hover.toggle"
      @mouseleave="hover.toggle"
      >Select</label
    >
    <input
      id="fileInput"
      class="hidden"
      type="file"
      @input.passive="handleInput"
    />
  </div>
</template>

<style scoped>
.dropzone {
  @apply flex flex-col items-center justify-center gap-y-4 rounded-2xl
         border-4 border-dashed border-gray-200 text-xl transition-all
         duration-75 ease-linear;
}
.input-button {
  @apply cursor-pointer rounded-2xl border-2 bg-green-300 px-5 py-2
         text-sm font-medium transition-all duration-75 ease-linear
         hover:drop-shadow-md;
}
</style>
