import { reactive } from "vue";

export const mainContent = reactive({
  current: "ResultTable",
  setCurrent(content: string) {
    this.current = content;
  },
});
