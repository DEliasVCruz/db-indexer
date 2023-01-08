package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc/search"
	"github.com/go-chi/chi/v5"
)

func SearchAdvance(w http.ResponseWriter, r *http.Request) {
	indexName := chi.URLParam(r, "indexName")

	var body search.ClientBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request body", http.StatusInternalServerError)
		return
	}

	if body.Pagination.From < 0 || body.Pagination.Size <= 0 {
		http.Error(w, "invalid page range", http.StatusBadRequest)
		return
	}

	var fieldsMust []*search.QueryType

	queryDataValues := reflect.ValueOf(*body.QueryData)
	queryDataTypes := queryDataValues.Type()

	for i := 0; i < queryDataValues.NumField(); i++ {

		if query := queryDataValues.Field(i).String(); query == "" {
			continue
		}

		fieldsMust = append(fieldsMust, &search.QueryType{
			Match: map[string]*search.Query{
				strings.ToLower(queryDataTypes.Field(i).Name): {
					Text: queryDataValues.Field(i).String(),
				},
			},
		})

	}

	bodyPayload := &search.SearchQuery{
		From: body.Pagination.From,
		Size: body.Pagination.Size,
		Query: &search.QueryType{
			Bool: &search.QueryBool{Must: fieldsMust},
		}}

	response, err := zinc.Search(indexName, bodyPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	searchResponse := data.BuildResponse(response)

	payLoad, err := json.Marshal(searchResponse)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	w.Write(payLoad)

}

func SearchField(w http.ResponseWriter, r *http.Request) {
	indexName := chi.URLParam(r, "indexName")

	queryParms := r.URL.Query()

	if err := check.ParamsOf("search", queryParms); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from, err := strconv.Atoi(queryParms.Get("from"))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	size, err := strconv.Atoi(queryParms.Get("size"))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if from < 0 || size <= 0 {
		http.Error(w, "invalid page range", http.StatusBadRequest)
		return
	}

	bodyPayload := &search.SearchQuery{
		From: from,
		Size: size,
		Query: &search.QueryType{
			Match: map[string]*search.Query{
				"contents": {Text: queryParms.Get("q")},
			},
		},
	}

	response, err := zinc.Search(indexName, bodyPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	searchResponse := data.BuildResponse(response)

	payLoad, err := json.Marshal(searchResponse)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	w.Write(payLoad)

}
