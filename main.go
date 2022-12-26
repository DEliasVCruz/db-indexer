package main

import (
	"os"
	"path/filepath"

	"github.com/DEliasVCruz/db-indexer/pkg/check"
	"github.com/DEliasVCruz/db-indexer/pkg/index"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
)

func main() {
	indexConfig, err := os.ReadFile("./index.json")
	check.Error("fileOpen", err)

	zinc.CreateIndex("appLogs", nil)

	indexer := index.Indexer{
		Name:       "emailsTest",
		DataFolder: filepath.Join(os.Args[1], "maildir/bailey-s/inbox"),
		Config:     indexConfig,
	}
	indexer.Index()
}
