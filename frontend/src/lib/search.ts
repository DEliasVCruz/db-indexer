import { request } from "@/lib/http";
import {
  Pagination,
  SearchObject,
  type QueryType,
  type SearchResponse,
} from "@/globals/types";

export async function search(
  searchType: string,
  searchQuery: QueryType,
  from: number,
  size: number,
  field: string
) {
  const url = new URL("http://localhost:3000/api/index/emailsTest/search");

  let response: Response;

  switch (searchType) {
    case "simple":
      if (typeof searchQuery.simple === "undefined") {
        return Promise.reject(new Error("empty simple query"));
      }
      response = await request.get({
        endpoint: url,
        urlParams: new URLSearchParams({
          q: searchQuery.simple,
          from: (from - 1).toString(),
          size: size.toString(),
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
        bodyPayload: new SearchObject(
          new Pagination(from - 1, size),
          searchQuery.advance
        ),
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
