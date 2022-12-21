package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var mainDir = filepath.Join(os.Args[1], "maildir")

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
