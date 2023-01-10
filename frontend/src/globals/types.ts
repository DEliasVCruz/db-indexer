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

export interface QueryData {
  from: string;
  to: string;
  subject: string;
  contents: string;
}

export interface Pager {
  from: number;
  size: number;
}

export interface AdvanceSearch {
  pagination: Pager;
  queryData: QueryData;
}

export interface QueryType {
  simple?: string;
  advance?: QueryData;
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
  lastQueryType: string;
  lastQuery: QueryType;
  nextPage(): void;
  prevPage(): void;
  resetRange(): void;
  setLastSimpleQuery(arg: string): void;
  setLastAdvanceQuery(arg: QueryData): void;
  setLastQueryType(arg: string): void;
  setTotalResults(arg: number): void;
  setEndRange(arg: number): void;
}

export class Pagination implements Pager {
  constructor(public from = 0, public size = 50) {
    this.from = from;
    this.size = size;
  }
}

export class MultiFieldQuery implements QueryData {
  constructor(
    public from = "",
    public to = "",
    public subject = "",
    public contents = ""
  ) {
    this.from = from;
    this.to = to;
    this.subject = subject;
    this.contents = contents;
  }
}

export class SearchObject implements AdvanceSearch {
  constructor(
    public pagination = new Pagination(),
    public queryData = new MultiFieldQuery()
  ) {
    this.pagination = pagination;
    this.queryData = queryData;
  }
}
