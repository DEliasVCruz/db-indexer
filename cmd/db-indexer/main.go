package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

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
		"content":                   "",
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
			metadataInfo = false
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
