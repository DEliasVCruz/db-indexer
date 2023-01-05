package search

type Query struct {
	From  int        `json:"from"`
	Size  int        `json:"size"`
	Query *QueryType `json:"query"`
}

type QueryType struct {
	Match map[string]map[string]string `json:"match"`
}
