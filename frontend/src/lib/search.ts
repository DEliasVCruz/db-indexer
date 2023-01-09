import { request } from "@/lib/http";
import type { QueryType, SearchResponse } from "@/globals/types";

export async function search(
  searchType: string,
  searchQuery: QueryType,
  from: string,
  size: string,
  field: string
) {
  const url = new URL("http://localhost:3000/index/emailsTest/search");

  let response: Response;

  switch (searchType) {
    case "simple":
      if (typeof searchQuery.simple === "undefined") {
        return Promise.reject(new Error("empty simple query"));
      }
      response = await request.get({
        endpoint: url,
        params: new URLSearchParams({
          q: searchQuery.simple,
          from: from,
          size: size,
          field: field,
        }),
      });
      break;
    case "advance":
      if (typeof searchQuery.advance === "undefined") {
        return Promise.reject(new Error("empty advance query object"));
      }
      response = await request.post({
        endpoint: url,
        body: searchQuery.advance,
      });
      break;
    default:
      return Promise.reject(
        new Error(`${searchType} is not a valid search type`)
      );
  }

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

  const total = data?.total;
  if (!total) {
    return Promise.reject(
      new Error("Error getting total number of found values")
    );
  }

  return { total: total, columns: columns };
}
