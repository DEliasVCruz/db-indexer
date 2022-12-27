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

	i.wg.Add(4)
	go i.findFiles(os.DirFS(i.DataFolder), files)
	go i.extractData(files, dataExtracts)
	go i.processData(dataExtracts, records)
	go i.collectRecords(records)

	i.wg.Wait()
}

func (i Indexer) findFiles(directory fs.FS, writeCh chan<- string) {
	defer i.wg.Done()

	fs.WalkDir(directory, ".", func(childPath string, dir fs.DirEntry, err error) error {

		fullPath := filepath.Join(i.DataFolder, childPath)
		if err != nil {
			zinc.LogError("appLogs", fmt.Sprintf("failed to read path %s", fullPath), err.Error())
			return nil
		}

		if !dir.IsDir() {
			writeCh <- fullPath
		}

		return nil
	})

	close(writeCh)
}

func (i Indexer) extractData(readCh <-chan string, writeCh chan<- map[string]string) {
	defer i.wg.Done()

	var wg sync.WaitGroup
	for file := range readCh {
		wg.Add(1)
		go data.Extract(file, writeCh, &wg)
	}

	wg.Wait()
	close(writeCh)

}

func (i Indexer) processData(readCh <-chan map[string]string, writeCh chan<- map[string]string) {
	defer i.wg.Done()

	var wg sync.WaitGroup
	for dataExtract := range readCh {
		wg.Add(1)
		go data.Process(dataExtract, writeCh, &wg)
	}

	wg.Wait()
	close(writeCh)

}

func (i Indexer) collectRecords(readCh <-chan map[string]string) {
	defer i.wg.Done()

	var records [100]map[string]string
	recordIdx := 0

	for record := range readCh {
		if recordIdx < 100 {
			records[recordIdx] = record
		} else {
			recordsSlice := make([]map[string]string, 100)
			copy(recordsSlice, records[:])

			i.wg.Add(1)
			go zinc.CreateDocBatch(i.Name, recordsSlice, i.wg)

			recordIdx = 0
			records[recordIdx] = record
		}

		recordIdx += 1
	}

	if recordIdx != 0 {
		recordsSlice := make([]map[string]string, recordIdx)
		copy(recordsSlice, records[:recordIdx])

		i.wg.Add(1)
		go zinc.CreateDocBatch(i.Name, recordsSlice, i.wg)
	}

}
