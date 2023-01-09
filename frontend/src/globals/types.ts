export interface ColumnData {
  name: string;
  values: string[];
}

export interface SearchResponse {
  data?: {
    columns: Array<ColumnData>;
    total: number;
  };
  error?: string;
}

export interface AdvanceSearch {
  pagination: {
    from: number;
    size: number;
  };
  queryData: {
    from: string;
    to: string;
    subject: string;
    contents: string;
  };
}

export interface Columns {
  columns: Array<ColumnData>;
  set(arg: Array<ColumnData>): void;
  getRow(arg: number): Map<string, string>;
}

export interface Results {
  total: number;
  size: number;
  from: number;
  to: number;
  page: number;
  lastQuery: string | AdvanceSearch;
  nextPage(): void;
  prevPage(): void;
  resetRange(): void;
  setLastQuery(arg: string): void;
  setTotalResults(arg: number): void;
  setEndRange(arg: number): void;
}

export class SearchObject implements AdvanceSearch {
  pagination: {
    from: number;
    size: number;
  };
  queryData: {
    from: string;
    to: string;
    subject: string;
    contents: string;
  };
  constructor() {
    this.pagination = {
      from: 0,
      size: 50,
    };
    this.queryData = {
      from: "",
      to: "",
      subject: "",
      contents: "",
    };
  }
}
