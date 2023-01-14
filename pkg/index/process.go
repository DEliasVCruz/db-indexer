package index

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
)

var fieldRegex = regexp.MustCompile(`^([\w\-]*):\s*(.*)`)
var brokenLineRegex = regexp.MustCompile(`^\s*(.*)\s*$`)
var messageRegex = regexp.MustCompile(`^<(\d+\.\d+)\..*`)

var fieldMetadataFlag = "x_filename"
var specialChars = [8]string{"-", "_", " ", "\n", ";", "=", `\`, "/"}

func (i Indexer) extract(data *data.DataInfo, ch chan<- map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	var input io.ReadCloser
	var err error
	var path string

	switch i.FileType {
	case "tar":
		path = data.TarBuf.Header.Name
		input, err = data.OpenTar()
	case "zip":
		path = data.RelPath
		input, err = i.dataFolder.Open(path)
	default:
		path, err = filepath.Abs(data.RelPath)
		if err != nil {
			log.Printf("failed to open file at relpath %s\n", data.RelPath)
		}
		input, err = os.Open(path)
	}

	check.Error("fileOpen", err)
	if err != nil {
		go zinc.LogError(
			"appLogs",
			fmt.Sprintf("failed to open file at path %s", path),
			err.Error(),
		)
		return
	}
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
		ch <- process(fields)
		return
	}

	log.Printf("data: failed to extract data at path %s", path)
	if scanner.Err() != nil {
		go zinc.LogError("appLogs", fmt.Sprintf("scanning failure at path %s", path), scanner.Err().Error())
		return
	}
	go zinc.LogError("appLogs", fmt.Sprintf("broken metadata at %s", path), "aborting indexing")
}

func process(fields map[string]string) map[string]string {

	messageId := messageRegex.FindStringSubmatch(fields["message_id"])
	if messageId != nil {
		fields["_id"] = messageId[1]
	}

	if val, ok := fields["content_type"]; ok {
		contentTypes := strings.Split(val, specialChars[4])
		if len(contentTypes) > 1 {
			fields["content_type"] = contentTypes[0]
			fields["charset"] = strings.Split(contentTypes[1], specialChars[5])[1]
		}
	}

	if val, ok := fields["x_folder"]; ok {
		fields["x_folder"] = strings.ReplaceAll(val, specialChars[6], specialChars[7])
	}

	return fields
}
