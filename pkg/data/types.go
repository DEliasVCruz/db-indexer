package data

type ColumnData struct {
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
