import { reactive } from "vue";
import {
  MultiFieldQuery,
  type ColumnData,
  type Columns,
  type Results,
} from "@/globals/types";
import { mainContent } from "@/globals/content";

export const column = reactive({
  selected: 0,
  select(column: number) {
    this.selected = column;
  },
});

export const results: Results = reactive({
  total: 50,
  size: 50,
  from: 1,
  to: 50,
  page: 1,
  lastQueryType: "",
  lastQuery: {
    simple: "",
    advance: new MultiFieldQuery(),
  },
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
  setLastSimpleQuery(newQuery) {
    this.lastQuery.simple = newQuery;
  },
  setLastAdvanceQuery(newQuery) {
    this.lastQuery.advance = newQuery;
  },
  setLastQueryType(queryType) {
    this.lastQueryType = queryType;
  },
  setTotalResults(newTotal) {
    this.total = newTotal;
  },
  setEndRange(newEnd) {
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
  columns: new Array<ColumnData>(),
  set(columns) {
    this.columns.length = 0;
    this.columns = columns;
  },
  getRow(rowId) {
    const rowData: Map<string, string> = new Map();
    this.columns.forEach((column) => {
      rowData.set(column.name, column.values[rowId]);
    });
    return rowData;
  },
});
