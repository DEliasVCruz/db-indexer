import { reactive } from "vue";

export const mainSection = reactive({
  current: "IndexExplorer",
  setCurrent(content: string) {
    this.current = content + "Explorer";
  },
});
