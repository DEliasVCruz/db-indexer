package index

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
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
	id         string
	archive    io.Reader
	wg         *sync.WaitGroup
}

func NewIndex(name, filetype, id string, upload *data.UploadData) {
	defer upload.File.Close()

	i := &Indexer{}

	i.Name = name
	i.wg = &sync.WaitGroup{}

	switch filetype {
	case "x-gzip":

		archive, err := gzip.NewReader(upload.File)

		if err != nil {

			log.Println(err.Error())
			if err := zinc.LogIndexStatus(true, "failed", i.id); err != nil {
				log.Println("Could not update index status")
				return
			}

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
			if err := zinc.LogIndexStatus(true, "failed", i.id); err != nil {
				log.Println("Could not update index status")
				return
			}

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

	if err := i.index(); err != nil {

		if err := zinc.LogIndexStatus(true, "failed", i.id); err != nil {
			log.Println("Could not update index status")
			return
		}

	}

	if err := zinc.LogIndexStatus(true, "done", i.id); err != nil {
		log.Println("Could not update index status")
	}
}

func (i Indexer) index() error {
	if !zinc.ExistsIndex(i.Name) {
		log.Printf("index: creating %s index", i.Name)
		if err := zinc.CreateIndex(i.Name, i.Config); err != nil {
			return err
		}
	}
	log.Printf("index: %s index already exists", i.Name)

	records := make(chan *data.Fields, 10)

	i.wg.Add(1)
	switch i.FileType {
	case "tar":
		go i.extractTAR(i.archive, records)
	case "fs":
		go i.extractFS(i.dataFolder, records)
	default:
		i.wg.Done()
		return errors.New(fmt.Sprintf("No matching extraction for filetype %s\n", i.FileType))
	}

	i.wg.Add(1)
	go i.collectRecords(records)

	i.wg.Wait()
	return nil
}

func (i Indexer) extractTAR(archive io.Reader, writeCh chan<- *data.Fields) {
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
					Err:    err,
					Size:   int(header.Size),
				},
				writeCh,
				&wg,
			)

		}

	}

	wg.Wait()
	close(writeCh)

}

func (i Indexer) extractFS(directory fs.FS, writeCh chan<- *data.Fields) {
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
			fileInfo, err := dir.Info()
			if err != nil {
				fmt.Printf("failed to get info for path %s", childPath)
			}

			wg.Add(1)
			go i.extract(
				&data.DataInfo{
					RelPath: childPath,
					Size:    int(fileInfo.Size()),
				}, writeCh, &wg)
		}

		return nil
	})

	wg.Wait()
	close(writeCh)
}

func (i Indexer) collectRecords(readCh <-chan *data.Fields) {
	defer i.wg.Done()

	var records [500]*data.Fields
	size := len(records)
	recordIdx := 0

	for record := range readCh {
		if recordIdx < size {
			records[recordIdx] = record
		} else {
			recordsSlice := make([]*data.Fields, size)
			copy(recordsSlice, records[:])

			i.wg.Add(1)
			go zinc.CreateDocBatch(i.Name, recordsSlice, i.wg)

			recordIdx = 0
			records[recordIdx] = record
		}

		recordIdx += 1
	}

	if recordIdx != 0 {
		recordsSlice := make([]*data.Fields, recordIdx)
		copy(recordsSlice, records[:recordIdx])

		i.wg.Add(1)
		go zinc.CreateDocBatch(i.Name, recordsSlice, i.wg)
	}

}
