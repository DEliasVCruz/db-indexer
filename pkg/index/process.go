package index

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/data"
)

var messageRegex = regexp.MustCompile(`^<(\d+\.\d+)\..*`)

func (i Indexer) extract(dataInfo *data.DataInfo, ch chan<- *data.Fields, wg *sync.WaitGroup) {
	defer wg.Done()

	var input io.ReadCloser
	var err error
	var path string

	switch i.FileType {
	case "tar":
		path = dataInfo.TarBuf.Header.Name
		input, err = dataInfo.OpenTar()
	case "fs":
		path = dataInfo.RelPath
		input, err = i.dataFolder.Open(path)
	default:
		log.Printf("failed to open filetype %v\n", i.FileType)
		return
	}

	check.Error("fileOpen", err)
	if err != nil {
		log.Printf("failed to open file %v", path)
		return
	}
	defer input.Close()

	var field string

	indexData := &data.Fields{}
	fieldName := &strings.Builder{}
	fieldValue := &strings.Builder{}

	fieldName.Grow(25)
	fieldValue.Grow(80)

	fileBuff := bufio.NewReaderSize(input, dataInfo.Size)
	fieldVals := reflect.ValueOf(indexData).Elem()

	remain := dataInfo.Size

	line, err := fileBuff.ReadSlice('\n')
	remain -= len(line)
	if err != nil && err != io.EOF {
		log.Printf("error reading %s file metadata", path)
		log.Println(err.Error())
		return
	}

	if err == io.EOF {
		log.Printf("error reading %s file metadata", path)
		log.Println(err.Error())
		return
	}

	sepIdx := bytes.IndexByte(line, ':')
	if sepIdx < 0 {
		log.Printf("invalid metadata at file %s", path)
		return
	}

	name := bytes.ReplaceAll(line[:sepIdx], []byte("-"), []byte(""))

	if len(name) > 25 {
		log.Printf("invalid metadata at file %s", path)
		return
	}

	if !bytes.Equal(name, []byte("MessageID")) {
		log.Printf("invalid metadata at file %s", path)
		return
	}

	fieldName.Write(name)
	field = fieldName.String()

	fieldValue.Write(bytes.TrimSpace(line[sepIdx+1:]))
	fieldVals.FieldByName(field).SetString(fieldValue.String())

	allMetadataParsed := false
	for !allMetadataParsed {
		line, err = fileBuff.ReadSlice('\n')
		remain -= len(line)
		if err != nil && err != io.EOF {
			log.Printf("error reading %s file metadata", path)
			log.Println(err.Error())
			return
		}

		if err == io.EOF {
			break
		}

		sepIdx = bytes.IndexByte(line, ':')
		if sepIdx < 0 {
			fieldValue.WriteString(" ")
			fieldValue.Write(bytes.TrimSpace(line))
			fieldVals.FieldByName(field).SetString(fieldValue.String())
			continue
		}

		name = bytes.ReplaceAll(line[:sepIdx], []byte("-"), []byte(""))

		if len(name) > 25 {
			fieldValue.WriteString(" ")
			fieldValue.Write(bytes.TrimSpace(line))
			fieldVals.FieldByName(field).SetString(fieldValue.String())
			continue
		}

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

		fieldName.Reset()
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
		return
	}

	fileBuff.ReadSlice('\n')
	fieldValue.WriteByte('\n')

	fieldValue.Grow(remain / 2)
	for {
		line, err := fileBuff.ReadSlice('\n')
		if err != nil && err != io.EOF {
			log.Printf("error reading %s file contents", path)
			return
		}

		fieldValue.Write(line)

		if err == io.EOF {
			break
		}

	}

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
