package zinc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DEliasVCruz/db-indexer/pkg/requests"
)

var request = requests.Request{
	BaseUrl:    "http://localhost",
	ServerPort: 4080,
	HttpClient: http.Client{Timeout: 30 * time.Second},
	Headers: map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("test_admin:test_password")),
		"Content-Type":  "application/json",
		"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
	},
	Retries: 3,
}

func ExistsIndex(index string) {
	status, _ := request.Get(fmt.Sprintf("api/index/%s", index), nil)
	if status == 200 {
		log.Printf("index: %s index already exists", index)
	} else {
		log.Printf("index: the %s index does not exist", index)
	}
}

func CreateIndex(index string, config []byte) {
	status, respBody := request.Post("api/index", config)

	if status == 200 {
		log.Printf("index: %s index was successfully created", index)
	} else {
		log.Fatalf("status: something went wrong got status code %d", status)
	}
	log.Printf("client: response body %s\n", respBody)
}

func DeleteIndex(index string) {
	status, _ := request.Delete(fmt.Sprintf("api/index/%s", index), nil)
	if status == 200 {
		log.Printf("index: %s index was deleted", index)
	} else {
		log.Fatalf("index: something went wrong trying to delete index %s", index)
	}
}

func CreateDoc(index string, payLoad map[string]string) {
	jsonPayLoad, _ := json.Marshal(payLoad)
	status, respBody := request.Post(fmt.Sprintf("api/%s/_doc", index), jsonPayLoad)
	if status == 200 {
		log.Printf("client: successful response with status %d and body %s", status, respBody)
	} else {
		log.Fatalf("client: could not index file with status %d and body %s", status, respBody)
	}
}

func CreateDocBatch(index string, payLoad []map[string]string) {
	jsonSlice, _ := json.Marshal(payLoad)
	jsonPayLoad := []byte(fmt.Sprintf(`{ "index": "%s", "records": %s }`, index, jsonSlice))
	status, respBody := request.Post("api/_bulkv2", jsonPayLoad)
	if status == 200 {
		log.Printf("client: successful response with status %d and body %s", status, respBody)
	} else {
		log.Fatalf("client: could not index file with status %d and body %s", status, respBody)
	}

}
