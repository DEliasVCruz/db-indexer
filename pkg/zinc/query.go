package zinc

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/search"
)

func Search(index string, searchQuery *search.SearchQuery) ([]byte, error) {

	jsonPayload, err := json.Marshal(searchQuery)
	if err != nil {
		log.Printf("server: internal server error %s", err)
		return nil, errors.New("could not marshal responsse")
	}

	status, body := request.Post("es/"+index+"/_search", jsonPayload)

	if err := check.SearchStatus(status); err != nil {
		return nil, err
	}

	return body, nil
}
