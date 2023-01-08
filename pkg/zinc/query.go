package zinc

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc/search"
)

func Search(index string, searchQuery *search.SearchQuery) (*search.Response, error) {

	jsonPayload, err := json.Marshal(searchQuery)
	if err != nil {
		log.Printf("server: internal server error %s", err)
		return nil, errors.New("could not marshal responsse")
	}

	status, body := request.Post("es/"+index+"/_search", jsonPayload)

	if err := check.SearchStatus(status); err != nil {
		return nil, err
	}

	var response *search.Response

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New("could not unmarshal responsse")
	}

	return response, nil
}
