package zinc

import (
	"encoding/json"
	"errors"
	"log"
)

type searchQuery struct {
	From  int        `json:"from"`
	Size  int        `json:"size"`
	Query *queryType `json:"query"`
}

type queryType struct {
	Match map[string]map[string]string `json:"match"`
}

func search(index string, payload *searchQuery) (int, []byte) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("server: internal server error %s", err)
	}

	return request.Post("es/"+index+"/_search", jsonPayload)

}

func SearchMatch(index string, query map[string]map[string]string, from, size int) ([]byte, error) {

	bodyPayload := &searchQuery{
		From: from,
		Size: size,
		Query: &queryType{
			Match: query,
		}}

	status, body := search(index, bodyPayload)

	if status != 200 {
		log.Printf("server: internal server error with status %d and body %s", status, body)
		return []byte(""), errors.New("index server error")
	}

	return body, nil
}
