package search

type Response struct {
	Took     int   `json:"took"`
	TimedOut bool  `json:"timed_out"`
	Hits     *hits `json:"hits"`
}

type hits struct {
	Total *total `json:"total"`
	Found []*hit `json:"hits"`
}

type total struct {
	Value int `json:"value"`
}

type hit struct {
	Index  string `json:"_index"`
	Id     string `json:"_id"`
	Source *Data  `json:"_source"`
}

type Data struct {
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
