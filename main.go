package main

import (
	"os"
	"path/filepath"

	"github.com/DEliasVCruz/db-indexer/pkg/index"
)

func main() {
	indexer := index.Indexer{
		Name:       "emailsTest",
		DataFolder: filepath.Join(os.Args[1], "maildir/bailey-s/inbox"),
	}
	indexer.Index()
}
