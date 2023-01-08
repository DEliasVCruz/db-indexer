<script setup lang="ts">
const props = defineProps<{
  field: string;
  modelValue: string;
}>();

const emit = defineEmits<{
  (event: "update:modelValue", payload: string): void;
}>();
</script>

<template>
  <div class="flex flex-row">
    <label :for="props.field" class="w-24 text-gray-500">{{
      props.field
    }}</label>
    <input
      :id="props.field"
      :value="modelValue"
      type="text"
      class="flex-grow border-b-2 border-b-gray-200 outline-none focus:border-b-blue-300"
      autocomplete="off"
      @input="
        emit('update:modelValue', ($event.target as HTMLInputElement).value)
      "
    />
  </div>
</template>
