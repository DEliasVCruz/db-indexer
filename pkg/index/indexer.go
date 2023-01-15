package index

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
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
	dataFolder fs.FS
	Config     []byte
	FileType   string
	archive    io.Reader
	wg         *sync.WaitGroup
}

func NewIndex(name, filetype string, upload *data.UploadData) {
	defer upload.File.Close()

	i := &Indexer{}

	i.Name = name
	i.wg = &sync.WaitGroup{}

	switch filetype {
	case "x-gzip":

		archive, err := gzip.NewReader(upload.File)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer archive.Close()

		i.archive = archive
		i.FileType = "tar"

	case "tar":

		i.archive = upload.File
		i.FileType = "tar"

	case "zip":

		zipFile, err := zip.NewReader(upload.File, upload.Size)
		if err != nil {
			log.Println(err.Error())
			return
		}
		i.dataFolder = zipFile
		i.FileType = "fs"

	case "folder":

		i.dataFolder = os.DirFS(upload.Folder)
		i.FileType = "fs"

	default:
		log.Printf("No matching indexer for filetype %s\n", filetype)
		return
	}

	i.index()
}

func (i Indexer) index() {
	if zinc.ExistsIndex(i.Name) == 200 {
		log.Printf("index: %s index already exists", i.Name)
	} else {
		log.Printf("index: the %s index does not exist", i.Name)
		log.Printf("index: creating %s index", i.Name)
		zinc.CreateIndex(i.Name, i.Config)
	}

	records := make(chan map[string]string)

	i.wg.Add(1)
	switch i.FileType {
	case "tar":
		go i.extractTAR(i.archive, records)
	case "fs":
		go i.extractFS(i.dataFolder, records)
	default:
		log.Printf("No matching extraction for filetype %s\n", i.FileType)
		i.wg.Done()
		return
	}

	i.wg.Add(1)
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
			buf := bytes.NewBuffer(make([]byte, 0, header.Size))
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
