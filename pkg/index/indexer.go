package index

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
)

type Indexer struct {
	Name       string
	DataFolder string
	Config     []byte
}

func (i Indexer) Index() {
	zinc.CreateIndex(i.Name, i.Config)
	i.indexDirFiles(os.DirFS(i.DataFolder))
}

func (i Indexer) indexDirFiles(directory fs.FS) {

	var records = []map[string]string{}

	fs.WalkDir(directory, ".", func(childPath string, dir fs.DirEntry, err error) error {

		fullPath := filepath.Join(i.DataFolder, childPath)
		if err != nil {
			log.Printf("error: when attempting to read file %s raised %s", fullPath, err)
			return nil
		}

		if !dir.IsDir() {
			fields, err := data.Extract(fullPath)
			if err != nil {
				log.Printf("error: %s", err)
				return nil
			}
			records = append(records, data.Process(fields))
		}

		if len(records) == 100 {
			zinc.CreateDocBatch(i.Name, records)
			records = nil
		}

		return nil
	})

	if len(records) != 0 {
		zinc.CreateDocBatch(i.Name, records)
		records = nil
	}
}
