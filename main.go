package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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

func mailDirIndex(mailDir fs.FS, rootPath string) {
	fs.WalkDir(mailDir, rootPath, indexWalker)
	if len(records) != 0 {
		createDocBatch(records)
		records = nil
	}
}

var records = []map[string]string{}

func indexWalker(childPath string, dir fs.DirEntry, err error) error {
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

	if len(records) == 100 {
		createDocBatch(records)
		records = nil
	}

	return nil
}

func main() {
	createIndex()
	mailDirIndex(os.DirFS(mainDir), "bailey-s/inbox")
}
