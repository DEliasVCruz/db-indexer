package zinc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DEliasVCruz/db-indexer/pkg/requests"
	"github.com/DEliasVCruz/db-indexer/pkg/search"
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
	status, _ := request.Get(fmt.Sprintf("api/index/%s", index), nil)

	return status
}

func CreateIndex(index string, config []byte) {
	status, body := request.Post("api/index", config)

	if status != 200 {
		log.Fatalf("index: could not create index, got status code %d and %s", status, body)
	}

	log.Printf("index: the %s index was created", index)
	LogInfo("appLogs", fmt.Sprintf("the %s idex was created", index))
}

func DeleteIndex(index string) {
	status, body := request.Delete(fmt.Sprintf("api/index/%s", index), nil)

	if status != 200 {
		log.Fatalf("index: something went wrong trying to delete index %s, got status code %d and %s", index, status, body)
	}
}

func CreateDoc(index string, payLoad []byte) int {
	status, body := request.Post(fmt.Sprintf("api/%s/_doc", index), payLoad)

	if status != 200 {
		log.Printf("client: could not index file with status %d", status)
		LogError("appLogs", fmt.Sprintf("could not index file"), string(body))
		return status
	}
	return status
}

func DeleteDoc(index, id string) {
	status, _ := request.Delete(fmt.Sprintf("api/%s/_doc/%s", index, id), nil)

	if status == 200 {
		log.Printf("client: deleted doc with id %s", id)
		LogInfo("appLogs", fmt.Sprintf("deleted doc with id %s", id))
	}

}

func CreateDocBatch(index string, payLoad []*search.Data, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonSlice, _ := json.Marshal(payLoad)
	jsonPayLoad := []byte(fmt.Sprintf(`{ "index": "%s", "records": %s }`, index, jsonSlice))
	status, body := request.Post("api/_bulkv2", jsonPayLoad)

	if status == 200 {
		log.Printf("client: %s", body)
		return
	}

	log.Printf("client: %s", body)
	for idx, record := range payLoad {
		body, _ = json.Marshal(record)
		if CreateDoc(index, body) != 200 && idx+1 != len(payLoad) {
			log.Printf("data: inserted %d records", idx)

			wg.Add(1)
			go CreateDocBatch(index, payLoad[idx+1:], wg)

			return
		}
	}

	fmt.Printf("data: inserted %d records", len(payLoad))

}
