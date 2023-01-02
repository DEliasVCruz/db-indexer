import { reactive } from "vue";

export const column = reactive({
  selected: "Contents",
  select(column: string) {
    this.selected = column;
  },
});

export const row = reactive({
  hovered: 0,
  hover(row: number) {
    this.hovered = row;
  },
});
