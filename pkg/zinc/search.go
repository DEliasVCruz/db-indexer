package zinc

import (
	"encoding/json"
	"errors"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc/search"
	"log"
)

func results(index string, payload *search.Query) (int, []byte) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("server: internal server error %s", err)
	}

	return request.Post("es/"+index+"/_search", jsonPayload)

}

func SearchMatch(index string, query map[string]map[string]string, from, size int) (*search.Response, error) {

	bodyPayload := &search.Query{
		From: from,
		Size: size,
		Query: &search.QueryType{
			Match: query,
		}}

	status, body := results(index, bodyPayload)

	if status != 200 {
		log.Printf("server: internal server error with status %d and body %s", status, body)
		return nil, errors.New("index server error")
	}

	var response *search.Response

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("could not unmarshal responsse")
	}

	return response, nil
}
