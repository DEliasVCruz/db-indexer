package main

import (
	"bufio"
	"bytes"
	"encoding/json"
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

	status := resp.StatusCode
	log.Printf("client: succsefull response with status code %d\n", status)

	body, err := ioutil.ReadAll(resp.Body)
	check("responseRead", err)

	return status, body
}

func createIndex() {
	status, _ := request(http.MethodHead, fmt.Sprintf("index/%s", defaultIndex), nil)
	if status == 200 {
		log.Printf("index: the %s index already exists", defaultIndex)
		status, _ := request(http.MethodDelete, fmt.Sprintf("index/%s", defaultIndex), nil)
		if status == 200 {
			log.Printf("index: the %s index was deleted", defaultIndex)
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
		log.Printf("index: %s was succsefully created", defaultIndex)
	} else {
		log.Printf("status: something went wrong got status code %d", status)
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

var field string
var data string

var mainDir = "enron_mail_20110402/maildir"

func fileIndexing(childPath string, dir fs.DirEntry, err error) error {
	fullPath := filepath.Join(mainDir, childPath)
	if err != nil {
		log.Printf("file: the following erro ocurred while attempting to read file %s - %s", fullPath, err)
		return fs.SkipDir
	}

	if !dir.IsDir() {
		input, err := os.Open(fullPath)
		check("fileOpen", err)
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
				} else {
					data = mailToRegex.FindStringSubmatch(line)[1]
					fields[field] += fmt.Sprintf(" %s", strings.TrimSpace(data))
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

		jsonPayLoad, _ := json.Marshal(fields)
		status, respBody := request(http.MethodPost, fmt.Sprintf("%s/_doc", defaultIndex), jsonPayLoad)
		if status == 200 {
			log.Printf("client: succsefull response with status code %d", status)
			log.Printf("client: response body %s", respBody)
		}

	}

	return nil
}

func main() {
	createIndex()
	fsys := os.DirFS(mainDir)
	fs.WalkDir(fsys, ".", fileIndexing)
}
