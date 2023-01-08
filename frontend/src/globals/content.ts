import { reactive } from "vue";

export const mainContent = reactive({
  current: "NoContent",
  setCurrent(content: string) {
    this.current = content;
  },
});
