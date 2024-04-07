import type { AdvanceSearch } from "@/globals/types";

interface RequestParams {
  endpoint: URL;
  urlParams?: URLSearchParams;
  bodyPayload?: AdvanceSearch;
  dataTransfer?: FormData;
}

export const request = {
  async get({ endpoint, urlParams: params }: RequestParams) {
    if (typeof params !== "undefined") {
      endpoint = new URL(`?${params.toString()}`, endpoint);
    }

    return fetch(endpoint);
  },

  async post({ endpoint, bodyPayload: body }: RequestParams) {
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

  async put({ endpoint, dataTransfer: data }: RequestParams) {
    if (typeof data === "undefined") {
      return Promise.reject(new Error("no request body was supplied"));
    }

    return fetch(endpoint, {
      method: "PUT",
      body: data,
    });
  },
};
