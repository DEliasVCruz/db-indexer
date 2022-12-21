package requests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
)

type Request struct {
	BaseUrl    string
	ServerPort int
	HttpClient http.Client
	Headers    map[string]string
}

func (r Request) Get(endpoint string, payLoad []byte) (int, []byte) {
	return baseRequest(r.HttpClient, http.MethodGet, fmt.Sprintf("%s:%d/%s", r.BaseUrl, r.ServerPort, endpoint), r.Headers, payLoad)
}

func (r Request) Post(endpoint string, payLoad []byte) (int, []byte) {
	return baseRequest(r.HttpClient, http.MethodPost, fmt.Sprintf("%s:%d/%s", r.BaseUrl, r.ServerPort, endpoint), r.Headers, payLoad)
}

func (r Request) Head(endpoint string, payLoad []byte) (int, []byte) {
	return baseRequest(r.HttpClient, http.MethodHead, fmt.Sprintf("%s:%d/%s", r.BaseUrl, r.ServerPort, endpoint), r.Headers, payLoad)
}

func (r Request) Delete(endpoint string, payLoad []byte) (int, []byte) {
	return baseRequest(r.HttpClient, http.MethodDelete, fmt.Sprintf("%s:%d/%s", r.BaseUrl, r.ServerPort, endpoint), r.Headers, payLoad)
}

func baseRequest(client http.Client, method, url string, headers map[string]string, payLoad []byte) (int, []byte) {
	req, err := http.NewRequest(method, url, bytes.NewReader(payLoad))
	check.Error("requestCreation", err)

	for header, value := range headers {
		req.Header.Set(header, value)
	}

	resp, err := client.Do(req)
	check.Error("requestAction", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check.Error("responseRead", err)

	return resp.StatusCode, body

}
