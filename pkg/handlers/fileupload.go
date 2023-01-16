package handlers

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/index"
	"github.com/DEliasVCruz/db-indexer/pkg/zinc"
	"github.com/go-chi/chi/v5"
)

const MAX_UPLOAD_SIZE = 500<<20 + 1024 // 500MB
const MAX_MEMORY = 20 << 20            // 20MB

func FileUpload(w http.ResponseWriter, r *http.Request) {

	indexName := chi.URLParam(r, "indexName")

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		http.Error(w, "file exceed max upload size", http.StatusBadRequest)
		return
	}

	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		file.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		file.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "application/zip" && filetype != "application/x-gzip" {

		tr := tar.NewReader(file)

		_, err := tr.Next()
		if err != nil {
			file.Close()
			http.Error(w, "the provided file format is not allowed", http.StatusBadRequest)
			return
		}

		filetype = "application/tar"
	}

	userid, err := r.Cookie("userid")
	if err != nil {
		uuid, err := exec.Command("uuidgen").Output()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userid := &http.Cookie{
			Name:    "userid",
			Value:   string(uuid),
			Expires: time.Now().Add(365 * 24 * time.Hour)}

		http.SetCookie(w, userid)
	}

	indexId, err := exec.Command("uuidgen").Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := string(indexId)

	go index.NewIndex(
		fmt.Sprintf("%s_%s", userid.String(), indexName),
		strings.TrimPrefix(filetype, "application/"),
		id,
		&data.UploadData{File: file, Size: fh.Size},
	)

	response, err := json.Marshal(
		&data.FileUploaded{
			Uploaded: true,
			State:    "processing",
			ID:       id,
		},
	)
	if err != nil {
		http.Error(w, "server marshaling failed", http.StatusInternalServerError)
		return
	}

	go zinc.CreateDoc("indexStatus", response)

	w.Write(response)

}
