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

func ExistsIndex(index string) int {
	status, body := request.Get(fmt.Sprintf("api/index/%s", index), nil)
	log.Printf("client: response with status %d and body %s\n", status, body)

	return status
}

func CreateIndex(index string, config []byte) {
	status, body := request.Post("api/index", config)
	log.Printf("client: response with status %d and body %s\n", status, body)

	if status == 200 {
		log.Printf("index: %s index was successfully created", index)
	} else {
		log.Fatalf("status: something went wrong got status code %d", status)
	}
}

func DeleteIndex(index string) {
	status, body := request.Delete(fmt.Sprintf("api/index/%s", index), nil)
	log.Printf("client: response with status %d and body %s\n", status, body)
	if status == 200 {
		log.Printf("index: %s index was deleted", index)
	} else {
		log.Fatalf("index: something went wrong trying to delete index %s", index)
	}
}

func CreateDoc(index string, payLoad map[string]string) {
	jsonPayLoad, _ := json.Marshal(payLoad)
	status, body := request.Post(fmt.Sprintf("api/%s/_doc", index), jsonPayLoad)
	log.Printf("client: response with status %d and body %s\n", status, body)
	if status == 200 {
		log.Printf("client: successful response with status %d and body %s", status, body)
	} else {
		log.Fatalf("client: could not index file with status %d and body %s", status, body)
	}
}

func CreateDocBatch(index string, payLoad []map[string][]byte) {
	jsonSlice, _ := json.Marshal(payLoad)
	jsonPayLoad := []byte(fmt.Sprintf(`{ "index": "%s", "records": %s }`, index, jsonSlice))
	status, body := request.Post("api/_bulkv2", jsonPayLoad)
	// This is repeating the message
	log.Printf("client: response with status %d and body %s\n", status, body)
	if status == 200 {
		log.Printf("client: successful response with status %d and body %s", status, body)
	} else {
		log.Fatalf("client: could not index file with status %d and body %s", status, body)
	}

}
