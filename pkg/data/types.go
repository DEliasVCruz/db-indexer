package data

import (
	"archive/tar"
	"bytes"
	"mime/multipart"
)

type Columns struct {
	From        []string `json:"from"`
	To          []string `json:"to"`
	Subject     []string `json:"subject"`
	Cc          []string `json:"cc"`
	Mime        []string `json:"mime_version"`
	ContentType []string `json:"content_type"`
	Charset     []string `json:"charset"`
	Encoding    []string `json:"content_transfer_encoding"`
	Bcc         []string `json:"bcc"`
	XFrom       []string `json:"x_from"`
	XTo         []string `json:"x_to"`
	Xcc         []string `json:"x_cc"`
	Xbcc        []string `json:"x_bcc"`
	XFolder     []string `json:"x_foler"`
	XOrigin     []string `json:"x_origin"`
	XFileName   []string `json:"x_filename"`
	FilPath     []string `json:"file_path"`
	Contents    []string `json:"contents"`
}

type SearchResponse struct {
	Data  *Data  `json:"data"`
	Error string `json:"error"`
}

type Data struct {
	Columns []*ColumnData `json:"columns"`
	Total   int           `json:"total"`
}

type ColumnData struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type DataInfo struct {
	RelPath string
	TarBuf  *TarBuf
	Size    int
	Err     error
}

type TarBuf struct {
	Buffer *bytes.Buffer
	Header *tar.Header
}

type UploadData struct {
	File   multipart.File
	Size   int64
	Folder string
}

type FileUploaded struct {
	Uploaded bool   `json:"uploaded,omitempty"`
	State    string `json:"state"`
	ID       string `json:"_id"`
}

type Fields struct {
	ID          string `json:"_id"`
	MessageID   string `json:"message_id"`
	Date        string `json:"date"`
	From        string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Cc          string `json:"cc"`
	Mime        string `json:"mime_version"`
	ContentType string `json:"content_type"`
	Charset     string `json:"charset"`
	Encoding    string `json:"content_transfer_encoding"`
	Bcc         string `json:"bcc"`
	XFrom       string `json:"x_from"`
	XTo         string `json:"x_to"`
	Xcc         string `json:"x_cc"`
	Xbcc        string `json:"x_bcc"`
	XFolder     string `json:"x_foler"`
	XOrigin     string `json:"x_origin"`
	XFileName   string `json:"x_filename"`
	FilePath    string `json:"file_path"`
	Contents    string `json:"contents"`
}
