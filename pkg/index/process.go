package index

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/search"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
)

var messageRegex = regexp.MustCompile(`^<(\d+\.\d+)\..*`)

func (i Indexer) extract(data *data.DataInfo, ch chan<- *search.Data, wg *sync.WaitGroup) {
	defer wg.Done()

	var input io.ReadCloser
	var err error
	var path string

	switch i.FileType {
	case "tar":
		path = data.TarBuf.Header.Name
		input, err = data.OpenTar()
	case "fs":
		path = data.RelPath
		input, err = i.dataFolder.Open(path)
	default:
		log.Printf("failed to open filetype %v\n", i.FileType)
		return
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

	var field string

	indexData := &search.Data{}
	fieldName := &strings.Builder{}
	fieldValue := &strings.Builder{}

	fieldName.Grow(120)
	fieldValue.Grow(120)

	fileBuff := bufio.NewReaderSize(input, data.Size)
	fieldVals := reflect.ValueOf(indexData).Elem()

	allMetadataParsed := false
	for !allMetadataParsed {
		line, err := fileBuff.ReadBytes('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("error reading %s file metadata", path)
			return
		}

		if err == io.EOF {
			break
		}

		sepIdx := bytes.IndexByte(line, ':')
		if sepIdx < 0 {
			fieldValue.WriteString(" ")
			fieldValue.Write(bytes.TrimSpace(line))
			fieldVals.FieldByName(field).SetString(fieldValue.String())
			continue
		}

		name := bytes.ReplaceAll(line[:sepIdx], []byte("-"), []byte(""))

		if bytes.Equal(name, []byte("XFileName")) {
			allMetadataParsed = true

			fieldValue.Reset()
			fieldValue.Write(bytes.TrimSpace(line[sepIdx+1:]))

			fieldVals.
				FieldByName("XFileName").
				SetString(fieldValue.String())

			fieldValue.Reset()

			break
		}

		fieldName.Write(name)

		if !fieldVals.FieldByName(fieldName.String()).IsValid() {
			fieldValue.WriteString(" ")
			fieldValue.Write(bytes.TrimSpace(line))
			fieldVals.FieldByName(field).SetString(fieldValue.String())
			continue
		}

		fieldValue.Reset()
		fieldValue.Write(bytes.TrimSpace(line[sepIdx+1:]))

		field = fieldName.String()
		fieldVals.
			FieldByName(field).
			SetString(fieldValue.String())

		fieldName.Reset()

	}

	if !allMetadataParsed {
		log.Printf("data: broken metadata at path %s", path)
		go zinc.LogError("appLogs", fmt.Sprintf("broken metadata at %s", path), "aborting indexing")
		return
	}

	fieldValue.Grow(fileBuff.Buffered())

	fileBuff.ReadBytes('\n')
	fieldValue.WriteByte('\n')

	fileBuff.WriteTo(fieldValue)

	fieldVals.FieldByName("Contents").SetString(fieldValue.String())
	fieldVals.FieldByName("FilePath").SetString(path)

	messageId := messageRegex.FindStringSubmatch(fieldVals.FieldByName("MessageID").String())
	if messageId != nil {
		fieldVals.FieldByName("ID").SetString(messageId[1])
	}

	if !fieldVals.FieldByName("ContentType").IsZero() {
		contentTypes := strings.Split(fieldVals.FieldByName("ContentType").String(), ";")
		if len(contentTypes) > 1 {
			fieldVals.FieldByName("ContentType").SetString(contentTypes[0])
			fieldVals.FieldByName("Charset").SetString(strings.Split(contentTypes[1], "=")[1])
		}
	}

	if !fieldVals.FieldByName("XFolder").IsZero() {
		xFolder := strings.ReplaceAll(fieldVals.FieldByName("XFolder").String(), `\`, "/")
		fieldVals.FieldByName("XFolder").SetString(xFolder)
	}

	ch <- indexData
}
