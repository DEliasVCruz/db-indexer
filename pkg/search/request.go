package search

type SearchQuery struct {
	From  int        `json:"from"`
	Size  int        `json:"size"`
	Query *QueryType `json:"query"`
}

type QueryType struct {
	Match map[string]*Query `json:"match,omitempty"`
	Bool  *QueryBool        `json:"bool,omitempty"`
	Ids   *Ids              `json:"ids,omitempty"`
}

type Query struct {
	Text string `json:"query"`
}

type QueryBool struct {
	Must []*QueryType `json:"must"`
}

type Ids struct {
	Values []string `json:"values"`
}
