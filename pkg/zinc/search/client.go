package search

type ClientBody struct {
	Pagination *Pager     `json:"pagination"`
	QueryData  *QueryData `json:"queryData"`
}

type QueryData struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Contents string `json:"contents"`
}

type Pager struct {
	From int `json:"from"`
	Size int `json:"size"`
}
