package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const serverPort = 4080

var client = http.Client{
	Timeout: 30 * time.Second,
}

func netwPost(endPoint string, payLoad *bytes.Reader) (int, []byte) {
	requestURL := fmt.Sprintf("http://localhost:%d/api/%s", serverPort, endPoint)
	req, err := http.NewRequest(http.MethodPost, requestURL, payLoad)
	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}
	req.SetBasicAuth("test_admin", "test_password")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client: error making http request: %s\n", err)
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	log.Printf("client: status code: %d\n", status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("client: could not read response body: %s\n", err)
	}

	return status, body
}

func createIndex() {
	jsonBody := []byte(`
						{
							"name": "testEmails",
							"storage_type": "disk",
							"shard_num": 3,
							"settings": {},
							"mappings": {
								"properties": {
									"message_id": {
										"type": "keyword",
										"index": true,
										"store": false
									},
									"date": {
										"type": "date",
										"format": "Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
										"time_zone": "",
										"index": true,
										"store": false,
										"sortable": true,
										"aggregatable": true
									},
									"from": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"to": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"subject": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"x_from": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"x_to": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"cc": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"x_bcc": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"x_folder": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									},
									"x_origin": {
										"type": "text",
										"index": true,
										"store": true,
										"sortable": true,
										"aggregatable": true
									}
								}
							}
						}
	`)
	jsonPayLoad := bytes.NewReader(jsonBody)
	status, respBody := netwPost("/index", jsonPayLoad)

	if status == 200 {
		fmt.Println("Evertthing is ok")
	} else {
		fmt.Printf("Something went wrogn got status code: %d", status)
	}
	fmt.Printf("client: response body: %s\n", respBody)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var field string
var data string

func main() {
	input, err := os.Open("./enron_mail_20110402/maildir/bailey-s/all_documents/10_")
	check(err)
	defer input.Close()

	var fields = map[string]string{
		"message_id":                "",
		"date":                      "",
		"from":                      "",
		"to":                        "",
		"subject":                   "",
		"cc":                        "",
		"mime_version":              "",
		"content_type":              "",
		"charset":                   "",
		"content_transfer_encoding": "",
		"bcc":                       "",
		"x_from":                    "",
		"x_to":                      "",
		"x_cc":                      "",
		"x_bcc":                     "",
		"x_folder":                  "",
		"x_origin":                  "",
		"x_filename":                "",
		"contents":                  "",
	}

	scanner := bufio.NewScanner(input)

	fieldRegex, _ := regexp.Compile(`^([\w\-]*): (.*)`)
	mailToRegex, _ := regexp.Compile(`^\s+(.*)\s*$`)
	regexMessage, _ := regexp.Compile(`^<(\d+\.\d+)\..*`)
	metadataInfo := true
	for scanner.Scan() {
		line := scanner.Text()
		if metadataInfo {
			if strings.TrimSpace(line) == "" {
				metadataInfo = false
			} else if fieldRegex.MatchString(line) {
				match := fieldRegex.FindStringSubmatch(line)
				field = strings.ReplaceAll(strings.ToLower(match[1]), "-", "_")
				data = strings.TrimSpace(match[2])
				fields[field] = strings.TrimSpace(data)
				fmt.Println(field)
				fmt.Println(data)
			} else {
				data = mailToRegex.FindStringSubmatch(line)[1]
				fields[field] += fmt.Sprintf(" %s", strings.TrimSpace(data))
				fmt.Println(data[1])
			}
		} else {
			fields["contents"] += fmt.Sprintf("%s\n", line)
		}
	}

	fields["message_id"] = regexMessage.FindStringSubmatch(fields["message_id"])[1]
	contentTypes := strings.Split(fields["content_type"], ";")
	fields["content_type"] = contentTypes[0]
	fields["charset"] = strings.Split(contentTypes[1], "=")[1]
	fields["x_folder"] = strings.ReplaceAll(fields["x_folder"], "\\", "/")

	myJson, _ := json.Marshal(fields)
	fmt.Println(string(myJson))
}
