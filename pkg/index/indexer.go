package index

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"sync"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
)

type Indexer struct {
	Name       string
	DataFolder string
	Config     []byte
	FileType   string
	ZipFolder  *zip.ReadCloser
	wg         *sync.WaitGroup
}

func (i Indexer) Index() {
	if zinc.ExistsIndex(i.Name) == 200 {
		log.Printf("index: %s index already exists", i.Name)
	} else {
		log.Printf("index: the %s index does not exist", i.Name)
		log.Printf("index: creating %s index", i.Name)
		zinc.CreateIndex(i.Name, i.Config)
	}

	records := make(chan map[string]string)

	i.wg = &sync.WaitGroup{}

	i.wg.Add(2)
	go i.extractFS(os.DirFS(i.DataFolder), records)
	go i.collectRecords(records)

	i.wg.Wait()
}

func (i Indexer) extractTAR(archive io.Reader, writeCh chan<- map[string]string) {
	tr := tar.NewReader(archive)

	var wg sync.WaitGroup
	done := false

	for {

		header, err := tr.Next()

		switch {
		case err == io.EOF:
			done = true
		case err != nil:
			go zinc.LogError(
				"appLogs",
				fmt.Sprintf("failed to read path %s", header.Name),
				err.Error(),
			)
			continue
		}

		if done {
			break
		}

		switch header.Typeflag {
		case tar.TypeDir:
			break
		case tar.TypeReg:
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(tr)
			go i.extract(
				&data.DataInfo{
					TarBuf: &data.TarBuf{Buffer: buf, Header: header},
					Err:    err},
				writeCh,
				&wg,
			)

		}

	}

	wg.Wait()
	close(writeCh)

}

func (i Indexer) extractFS(directory fs.FS, writeCh chan<- map[string]string) {
	defer i.wg.Done()

	var wg sync.WaitGroup
	fs.WalkDir(directory, ".", func(childPath string, dir fs.DirEntry, err error) error {

		if err != nil {
			go zinc.LogError(
				"appLogs",
				fmt.Sprintf("failed to read path %s", childPath),
				err.Error(),
			)
			return nil
		}

		if !dir.IsDir() {
			wg.Add(1)
			go i.extract(&data.DataInfo{RelPath: childPath}, writeCh, &wg)
		}

		return nil
	})

	wg.Wait()
	close(writeCh)
}

func (i Indexer) collectRecords(readCh <-chan map[string]string) {
	defer i.wg.Done()

	var records [500]map[string]string
	size := len(records)
	recordIdx := 0

	for record := range readCh {
		if recordIdx < size {
			records[recordIdx] = record
		} else {
			recordsSlice := make([]map[string]string, size)
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
