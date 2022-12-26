package data

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
)

var fieldRegex, _ = regexp.Compile(`^([\w\-]*):\s*(.*)`)
var brokenLineRegex, _ = regexp.Compile(`^\s*(.*)\s*$`)
var messageRegex, _ = regexp.Compile(`^<(\d+\.\d+)\..*`)

var fieldMetadataFlag = "x_filename"
var specialChars = [8]string{"-", "_", " ", "\n", ";", "=", `\`, "/"}

func Extract(path string) (map[string]string, error) {
	log.Printf("file: reading file path %s", path)
	input, err := os.Open(path)
	check.Error("fileOpen", err)
	defer input.Close()

	fields := make(map[string]string)
	var field string

	allMetadataParsed := false
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if !allMetadataParsed {
			if field == fieldMetadataFlag {
				allMetadataParsed = true
			} else if match := fieldRegex.FindStringSubmatch(line); match != nil {
				field = strings.ReplaceAll(strings.ToLower(match[1]), specialChars[0], specialChars[1])
				fields[string(field)] = strings.TrimSpace(match[2])
			} else {
				fields[field] += specialChars[2] + strings.TrimSpace(brokenLineRegex.FindStringSubmatch(line)[1])
			}
		} else {
			fields["contents"] += line + specialChars[3]
		}
	}

	if !allMetadataParsed {
		return fields, errors.New(fmt.Sprintf("broken metadata at %s aborting indexing", path))
	}

	return fields, nil
}

func Process(fields map[string]string) map[string]string {
	messageId := messageRegex.FindStringSubmatch(fields["message_id"])
	if messageId != nil {
		fields["message_id"] = messageId[1]
	}
	contentTypes := strings.Split(fields["content_type"], specialChars[4])
	if len(contentTypes) > 1 {
		fields["content_type"] = contentTypes[0]
		fields["charset"] = strings.Split(contentTypes[1], specialChars[5])[1]
	}
	fields["x_folder"] = strings.ReplaceAll(fields["x_folder"], specialChars[6], specialChars[7])

	return fields
}
