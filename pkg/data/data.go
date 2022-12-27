package data

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
)

var fieldRegex = regexp.MustCompile(`^([\w\-]*):\s*(.*)`)
var brokenLineRegex = regexp.MustCompile(`^\s*(.*)\s*$`)
var messageRegex = regexp.MustCompile(`^<(\d+\.\d+)\..*`)

var fieldMetadataFlag = "x_filename"
var specialChars = [8]string{"-", "_", " ", "\n", ";", "=", `\`, "/"}

func Extract(path string, ch chan<- map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	input, err := os.Open(path)
	check.Error("fileOpen", err)
	defer input.Close()

	fields := make(map[string]string)
	var field string

	scanner := bufio.NewScanner(input)
	buf := make([]byte, 70*1024)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)

	allMetadataParsed := false
	for scanner.Scan() {
		line := scanner.Text()
		if !allMetadataParsed {
			if match := fieldRegex.FindStringSubmatch(line); match != nil {
				field = strings.ReplaceAll(strings.ToLower(match[1]), specialChars[0], specialChars[1])
				fields[field] = strings.TrimSpace(match[2])
				if field == fieldMetadataFlag {
					allMetadataParsed = true
				}
			} else {
				fields[field] += specialChars[2] + strings.TrimSpace(brokenLineRegex.FindStringSubmatch(line)[1])
			}
		} else {
			fields["contents"] += line + specialChars[3]
		}
	}

	fields["file_path"] = path

	if allMetadataParsed {
		ch <- fields
		return
	}

	zinc.LogError("appLogs", fmt.Sprintf("broken metadata at %s", path), "aborting indexing")
}

func Process(fields map[string]string, ch chan<- map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	messageId := messageRegex.FindStringSubmatch(fields["message_id"])
	if messageId != nil {
		fields["_id"] = messageId[1]
	}
	contentTypes := strings.Split(fields["content_type"], specialChars[4])
	if len(contentTypes) > 1 {
		fields["content_type"] = contentTypes[0]
		fields["charset"] = strings.Split(contentTypes[1], specialChars[5])[1]
	}
	fields["x_folder"] = strings.ReplaceAll(fields["x_folder"], specialChars[6], specialChars[7])

	ch <- fields
}
