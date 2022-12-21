package data

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
)

var fieldRegex, _ = regexp.Compile(`^([\w\-]*): (.*)`)
var brokenLineRegex, _ = regexp.Compile(`^\s*(.*)\s*$`)
var messageRegex, _ = regexp.Compile(`^<(\d+\.\d+)\..*`)

func Extract(path string) (map[string]string, error) {
	log.Printf("file: reading file path %s", path)
	input, err := os.Open(path)
	check.Error("fileOpen", err)
	defer input.Close()

	fields := map[string]string{}
	field := []byte("")

	allMetadataParsed := false
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !allMetadataParsed {
			if bytes.Equal(field, []byte("x_filename")) {
				allMetadataParsed = true
			} else if fieldRegex.Match(line) {
				match := fieldRegex.FindSubmatch(line)
				field = bytes.ReplaceAll(bytes.ToLower(match[1]), []byte("-"), []byte("_"))
				fields[string(field)] = string(bytes.TrimSpace(match[2]))
			} else {
				fields[string(field)] += fmt.Sprintf(" %s", bytes.TrimSpace(brokenLineRegex.FindSubmatch(line)[1]))
			}
		} else {
			fields["contents"] += fmt.Sprintf("%s\n", line)
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
	contentTypes := strings.Split(fields["content_type"], ";")
	if len(contentTypes) > 1 {
		fields["content_type"] = contentTypes[0]
		fields["charset"] = strings.Split(contentTypes[1], "=")[1]
	}
	fields["x_folder"] = strings.ReplaceAll(fields["x_folder"], "\\", "/")

	return fields
}
