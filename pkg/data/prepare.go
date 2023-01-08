package data

import (
	"reflect"

	"github.com/DEliasVCruz/db-indexer/pkg/zinc/search"
)

func BuildResponse(response *search.Response) *SearchResponse {
	columns := &Columns{}

	columnsValues := reflect.ValueOf(columns).Elem()

	for _, hit := range response.Hits.Found {
		hitValues := reflect.ValueOf(*hit.Source)
		for i := 0; i < hitValues.NumField(); i++ {
			columnsValues.Field(i).Set(reflect.Append(columnsValues.Field(i), hitValues.Field(i)))
		}
	}

	columnsData := []*ColumnData{}

	columnsValues = reflect.ValueOf(*columns)
	columnsTypes := columnsValues.Type()
	for i := 0; i < columnsValues.NumField(); i++ {
		columnsData = append(columnsData, &ColumnData{
			Name:   columnsTypes.Field(i).Name,
			Values: columnsValues.Field(i).Interface().([]string),
		})
	}

	return &SearchResponse{Data: &Data{Columns: columnsData, Total: response.Hits.Total.Value}}
}
