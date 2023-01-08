import { reactive } from "vue";
import type { ColumnData } from "@/globals/types";
import { mainContent } from "@/globals/content";

interface Columns {
  columns: Array<ColumnData>;
  set(columns: Array<ColumnData>): void;
  getRow(row: number): Map<string, string>;
}

export const column = reactive({
  selected: "Contents",
  select(column: string) {
    this.selected = column;
  },
});

export const results = reactive({
  total: 50,
  size: 50,
  from: 1,
  to: 50,
  page: 1,
  lastQuery: "",
  nextPage() {
    this.page += 1;
    this.from = this.size * (this.page - 1) + 1;

    const upperRange = this.size * this.page;
    if (upperRange > this.total) {
      this.to = this.total;
    } else {
      this.to = upperRange;
    }
  },
  prevPage() {
    this.page -= 1;

    const gap = this.to - (this.from - 1);
    if (this.to < this.total) {
      this.to -= this.size;
    } else if (gap < this.size) {
      this.to -= gap;
    } else {
      this.to -= this.size;
    }

    this.from -= this.size;
  },
  resetRange() {
    this.from = 1;
    if (this.total < this.size) {
      this.to = this.total;
    } else {
      this.to = this.size;
    }
    this.page = 1;
  },
  setLastQuery(newQuery: string) {
    this.lastQuery = newQuery;
  },
  setTotalResults(newTotal: number) {
    this.total = newTotal;
  },
  setEndRange(newEnd: number) {
    this.to = newEnd;
  },
});

export const row = reactive({
  hovered: 0,
  doubleClicked: 0,
  data: new Map([["", ""]]),
  hover(row: number) {
    this.hovered = row;
  },
  render(rowId: number) {
    this.data = columnData.getRow(rowId);
    mainContent.setCurrent("MailView");
  },
});

export const columnData: Columns = reactive({
  columns: new Array(),
  set(columns: Array<ColumnData>) {
    this.columns.length = 0;
    this.columns = columns;
  },
  getRow(rowId: number) {
    const rowData: Map<string, string> = new Map();
    this.columns.forEach((column) => {
      rowData.set(column.name, column.values[rowId]);
    });
    return rowData;
  },
});
