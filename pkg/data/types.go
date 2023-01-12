package data

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
	XCc         []string `json:"x_cc"`
	XBcc        []string `json:"x_bcc"`
	XFolder     []string `json:"x_foler"`
	XOrigin     []string `json:"x_origin"`
	XFilename   []string `json:"x_filename"`
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
	Err     error
}

type TarBuf struct {
	Buffer *bytes.Buffer
	Header *tar.Header
}
