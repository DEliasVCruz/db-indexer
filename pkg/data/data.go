package data

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
)

var fieldRegex, _ = regexp.Compile(`^([\w\-]*): (.*)`)
var brokenLineRegex, _ = regexp.Compile(`^\s*(.*)\s*$`)
var messageRegex, _ = regexp.Compile(`^<(\d+\.\d+)\..*`)

var fieldBreakBefore, fieldBreakAfter = []byte("-"), []byte("_")
var fieldMetadataFlag, emptyByte = []byte("x_filename"), []byte("")
var dataJoin, dataNewLine = []byte(" "), []byte("\n")
var contentTypeSep, charsetSep = []byte(";"), []byte("=")
var oldPathSep, newPathSep = []byte(`\`), []byte("/")

func Extract(path string) (map[string][]byte, error) {
	log.Printf("file: reading file path %s", path)
	input, err := os.Open(path)
	check.Error("fileOpen", err)
	defer input.Close()

	fields := make(map[string][]byte)
	var field []byte

	allMetadataParsed := false
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !allMetadataParsed {
			if bytes.Equal(field, fieldMetadataFlag) {
				allMetadataParsed = true
			} else if match := fieldRegex.FindSubmatch(line); match != nil {
				field = bytes.ReplaceAll(bytes.ToLower(match[1]), fieldBreakBefore, fieldBreakAfter)
				fields[string(field)] = bytes.TrimSpace(match[2])
			} else {
				byteData := [][]byte{fields[string(field)], bytes.TrimSpace(brokenLineRegex.FindSubmatch(line)[1])}
				fields[string(field)] = bytes.Join(byteData, dataJoin)
			}
		} else {
			byteData := [][]byte{fields["contents"], line, dataNewLine}
			fields["contents"] = bytes.Join(byteData, emptyByte)
		}
	}

	if !allMetadataParsed {
		return fields, errors.New(fmt.Sprintf("broken metadata at %s aborting indexing", path))
	}

	return fields, nil
}

func Process(fields map[string][]byte) map[string][]byte {
	messageId := messageRegex.FindSubmatch(fields["message_id"])
	if messageId != nil {
		fields["message_id"] = messageId[1]
	}
	contentTypes := bytes.Split(fields["content_type"], contentTypeSep)
	if len(contentTypes) > 1 {
		fields["content_type"] = contentTypes[0]
		fields["charset"] = bytes.Split(contentTypes[1], charsetSep)[1]
	}
	fields["x_folder"] = bytes.ReplaceAll(fields["x_folder"], oldPathSep, newPathSep)

	return fields
}
