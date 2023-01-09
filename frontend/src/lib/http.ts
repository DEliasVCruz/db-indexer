import type { AdvanceSearch } from "@/globals/types";

interface GetParams {
  endpoint: URL;
  params?: URLSearchParams;
  body?: AdvanceSearch;
}

export const request = {
  async get({ endpoint, params }: GetParams) {
    let url: URL;

    if (typeof params !== "undefined") {
      url = new URL(`?${params.toString()}`, endpoint);
    } else {
      url = endpoint;
    }

    return fetch(url);
  },
  async post({ endpoint, body }: GetParams) {
    if (typeof body === "undefined") {
      return Promise.reject(new Error("no request body was supplied"));
    }

    return fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json;charset=utf-8",
      },
      body: JSON.stringify(body),
    });
  },
};
