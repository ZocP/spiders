package abstract

type PairQA struct {
	Q string
	A string
}

type ArticleQA struct {
	Link  string
	Title string
	QA    []PairQA
	Mark  string
}
