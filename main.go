package main

import (
	"os"
	"path/filepath"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/index"
)

func main() {
	indexConfig, err := os.ReadFile("./index.json")
	check.Error("fileOpen", err)

	indexer := index.Indexer{
		Name:       "emailsTest",
		DataFolder: os.Args[1],
		Config:     indexConfig,
	}
	indexer.Index()
}
