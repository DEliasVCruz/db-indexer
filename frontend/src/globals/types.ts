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
