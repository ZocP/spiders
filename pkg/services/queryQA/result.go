package queryQA

type Result struct {
	Title string   `json:"title"`
	Link  string   `json:"link"`
	QA    []PairQA `json:"qa"`
}

type PairQA struct {
	Q string `json:"q"`
	A string `json:"a"`
}
