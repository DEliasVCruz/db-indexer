package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const serverPort = 4080
const defaultIndex = "emails"

var client = http.Client{
	Timeout: 30 * time.Second,
}

func request(method, endPoint string, payLoad []byte) (int, []byte) {
	bodyReader := bytes.NewReader(payLoad)
	requestURL := fmt.Sprintf("http://localhost:%d/api/%s", serverPort, endPoint)

	req, err := http.NewRequest(method, requestURL, bodyReader)
	check("requestCreation", err)

	req.SetBasicAuth("test_admin", "test_password")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := client.Do(req)
	check("requestAction", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check("responseRead", err)

	return resp.StatusCode, body
}

func createIndex() {
	status, _ := request(http.MethodHead, fmt.Sprintf("index/%s", defaultIndex), nil)
	if status == 200 {
		log.Printf("index: %s index already exists", defaultIndex)
		status, _ := request(http.MethodDelete, fmt.Sprintf("index/%s", defaultIndex), nil)
		if status == 200 {
			log.Printf("index: %s index was deleted", defaultIndex)
		} else {
			log.Fatalf("index: something went wrong trying to delete index %s", defaultIndex)
		}
	} else {
		log.Printf("index: the %s index does not exist", defaultIndex)
	}
	jsonBody, err := os.ReadFile("./index.json")
	check("fileOpen", err)

	status, respBody := request(http.MethodPost, "index", jsonBody)

	if status == 200 {
		log.Printf("index: %s index was successfully created", defaultIndex)
	} else {
		log.Fatalf("status: something went wrong got status code %d", status)
	}
	log.Printf("client: response body %s\n", respBody)
}

func check(check string, err error) {
	if err != nil {
		switch check {
		case "requestCreation":
			log.Fatalf("client: could not create request with %s\n", err)
		case "fileOpen":
			log.Fatalf("file: could not open given file with %s\n", err)
		case "requestAction":
			log.Fatalf("client: error making http request with %s\n", err)
		case "responseRead":
			log.Fatalf("client: could not read response body with %s\n", err)
		default:
			log.Fatalf("error: something went wrong with %s\n", err)
		}
	}
}

var mainDir = filepath.Join(os.Args[1], "maildir")

var fieldRegex, _ = regexp.Compile(`^([\w\-]*): (.*)`)
var brokenLineRegex, _ = regexp.Compile(`^\s*(.*)\s*$`)
var messageRegex, _ = regexp.Compile(`^<(\d+\.\d+)\..*`)

func dataExtract(path string) (map[string]string, error) {
	log.Printf("file: reading file path %s", path)
	input, err := os.Open(path)
	check("fileOpen", err)
	defer input.Close()

	fields := map[string]string{}
	field := ""
	data := ""

	allMetadataParsed := false
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if !allMetadataParsed {
			if field == "x_filename" {
				allMetadataParsed = true
			} else if fieldRegex.MatchString(line) {
				match := fieldRegex.FindStringSubmatch(line)
				field = strings.ReplaceAll(strings.ToLower(match[1]), "-", "_")
				data = strings.TrimSpace(match[2])
				fields[field] = strings.TrimSpace(data)
			} else {
				data = brokenLineRegex.FindStringSubmatch(line)[1]
				fields[field] += fmt.Sprintf(" %s", strings.TrimSpace(data))
			}
		} else {
			fields["contents"] += fmt.Sprintf("%s\n", line)
		}
	}

	if !allMetadataParsed {
		return fields, errors.New(fmt.Sprintf("broken metadata at %s aborting indexing", path))
	}

	messageId := messageRegex.FindStringSubmatch(fields["message_id"])
	if messageId != nil {
		fields["message_id"] = messageId[1]
	}
	contentTypes := strings.Split(fields["content_type"], ";")
	if len(contentTypes) > 1 {
		fields["content_type"] = contentTypes[0]
		fields["charset"] = strings.Split(contentTypes[1], "=")[1]
	}
	fields["x_folder"] = strings.ReplaceAll(fields["x_folder"], "\\", "/")

	return fields, nil

}

func createDocBatch(payLoad []map[string]string) {
	jsonSlice, _ := json.Marshal(payLoad)
	jsonPayLoad := []byte(fmt.Sprintf(`{ "index": "%s", "records": %s }`, defaultIndex, jsonSlice))
	status, respBody := request(http.MethodPost, "_bulkv2", jsonPayLoad)
	if status == 200 {
		log.Printf("client: successful response with status %d and body %s", status, respBody)
	} else {
		log.Fatalf("client: could not index file with status %d and body %s", status, respBody)
	}

}

func createDoc(payLoad map[string]string) {
	jsonPayLoad, _ := json.Marshal(payLoad)
	status, respBody := request(http.MethodPost, fmt.Sprintf("%s/_doc", defaultIndex), jsonPayLoad)
	if status == 200 {
		log.Printf("client: successful response with status %d and body %s", status, respBody)
	} else {
		log.Fatalf("client: could not index file with status %d and body %s", status, respBody)
	}
}

var records = []map[string]string{}

func fsWalker(childPath string, dir fs.DirEntry, err error) error {
	fullPath := filepath.Join(mainDir, childPath)
	if err != nil {
		log.Printf("error: when attempting to read file %s raised %s", fullPath, err)
		return nil
	}

	if !dir.IsDir() {
		fields, err := dataExtract(fullPath)
		if err != nil {
			log.Printf("error: %s", err)
			return nil
		}
		records = append(records, fields)
	}

	if len(records) == 2 {
		createDocBatch(records)
		records = nil
	}

	return nil
}

func main() {
	createIndex()
	fsys := os.DirFS(mainDir)
	fs.WalkDir(fsys, ".", fsWalker)
	if len(records) != 0 {
		createDocBatch(records)
		records = nil
	}
}
