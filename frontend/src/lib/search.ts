import { request } from "@/lib/http";
import type { SearchResponse } from "@/globals/types";

export async function searchText(text: string) {
  const url = new URL("http://localhost:3001/index/myindex/search");

  const response = await request
    .get({
      endpoint: url,
      params: new URLSearchParams({ q: text }),
    })
    .catch((error: Error) => {
      return Promise.reject(error);
    });

  const { data, error }: SearchResponse = await response.json();
  if (!response.ok) {
    return Promise.reject(new Error(`An error has ocurred and is ${error}`));
  }

  const columns = data?.columns;
  if (!columns) {
    return Promise.reject(
      new Error("No match found for given query, please try a new one")
    );
  }

  return columns;
}
