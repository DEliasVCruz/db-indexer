import { reactive } from "vue";
import type { ColumnData } from "@/globals/types";

interface Columns {
  columns: Array<ColumnData>;
  set(columns: Array<ColumnData>): void;
}

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

export const columnData: Columns = reactive({
  columns: new Array(),
  set(columns: Array<ColumnData>) {
    this.columns.length = 0;
    this.columns = columns;
  },
});
