package index

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
)

type Indexer struct {
	Name       string
	DataFolder string
	Config     []byte
	wg         *sync.WaitGroup
	records    [100]map[string]string
	recordIdx  int
}

func (i Indexer) Index() {
	if zinc.ExistsIndex(i.Name) == 200 {
		log.Printf("index: %s index already exists", i.Name)
	} else {
		log.Printf("index: the %s index does not exist", i.Name)
		zinc.CreateIndex(i.Name, i.Config)
	}

	files := make(chan string)
	records := make(chan map[string]string)
	dataExtracts := make(chan map[string]string)

	i.wg = &sync.WaitGroup{}

	i.wg.Add(1)
	go i.findFiles(os.DirFS(i.DataFolder), files)

	go i.extractData(files, dataExtracts)

	go i.processData(dataExtracts, records)

	go i.collectRecords(records)

	i.wg.Wait()

	close(dataExtracts)
	close(records)
	if i.recordCount != 0 {
		zinc.CreateDocBatch(i.Name, i.records[:i.recordCount])
	}

}

func (i Indexer) findFiles(directory fs.FS, ch chan<- string) {
	defer i.wg.Done()

	fs.WalkDir(directory, ".", func(childPath string, dir fs.DirEntry, err error) error {

		fullPath := filepath.Join(i.DataFolder, childPath)
		if err != nil {
			zinc.LogError(fmt.Sprintf("failed to read path %s", fullPath), err.Error())
			return nil
		}

		if !dir.IsDir() {
			ch <- fullPath
		}

		return nil
	})

	close(ch)
}

func (i Indexer) extractData(readCh <-chan string, writeCh chan<- map[string]string) {
	for file := range readCh {
		i.wg.Add(1)
		go data.Extract(file, writeCh, i.wg)
	}
}

func (i Indexer) processData(readCh <-chan map[string]string, writeCh chan<- map[string]string) {
	for dataExtract := range readCh {
		i.wg.Add(1)
		go data.Process(dataExtract, writeCh, i.wg)
	}

}

func (i Indexer) collectRecords(readCh <-chan map[string]string) {
	for record := range readCh {
		i.wg.Add(1)
		if i.recordIdx < 100 {
			i.records[i.recordIdx] = record
			i.recordIdx += 1
		} else {
			zinc.CreateDocBatch(i.Name, i.records[:i.recordIdx])
			i.recordIdx = 0
		}
		i.wg.Done()
	}

}
