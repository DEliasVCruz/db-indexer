package data

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/DEliasVCruz/db-indexer/pkg/search"
)

func BuildSearchResponse(responseBody []byte) (*SearchResponse, error) {
	var response *search.Response

	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, errors.New("could not unmarshal responsse")
	}

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

	return &SearchResponse{Data: &Data{
		Columns: columnsData,
		Total:   response.Hits.Total.Value}}, nil
}

func BuildIndexSearchResponse(responseBody []byte) (*FileUploaded, error) {
	var response *search.IndexStatusResponse

	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, errors.New("could not unmarshal responsse")
	}

	indexStatus := response.Hits.Found[0].Source

	return &FileUploaded{
		Uploaded: indexStatus.Uploaded,
		State:    indexStatus.State,
		ID:       indexStatus.ID,
	}, nil

}
