interface GetParams {
  endpoint: URL;
  params?: URLSearchParams;
}

export const request = {
  async get({ endpoint, params }: GetParams) {
    let url: URL;

    if (typeof params !== "undefined") {
      url = new URL(`?${params.toString()}`, endpoint);
    } else {
      url = endpoint;
    }

    return await fetch(url);
  },
};
