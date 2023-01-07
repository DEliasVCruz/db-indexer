import { reactive } from "vue";

export const mainContent = reactive({
  current: "resultTable",
  setCurrent(content: string) {
    this.current = content;
  },
});
