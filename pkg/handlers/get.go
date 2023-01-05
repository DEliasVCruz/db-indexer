package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
	"github.com/go-chi/chi/v5"
)

func SearchContents(w http.ResponseWriter, r *http.Request) {
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

	if from >= size || from < 0 || size <= 0 {
		http.Error(w, "invalid page range", http.StatusBadRequest)
		return
	}

	query := map[string]map[string]string{
		"contents": {
			"query": queryParms.Get("q"),
		},
	}

	response, err := zinc.SearchMatch(indexName, query, from, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	data := &data.ColumnData{}

	dataValues := reflect.ValueOf(data).Elem()

	for _, hit := range response.Hits.Found {
		hitValues := reflect.ValueOf(*hit.Source)
		for i := 0; i < hitValues.NumField(); i++ {
			dataValues.Field(i).Set(reflect.Append(dataValues.Field(i), hitValues.Field(i)))
		}
	}

	payLoad, err := json.Marshal(data)

	w.Write(payLoad)

}
