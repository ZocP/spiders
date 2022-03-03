package qa

import "qa_spider/pkg/spiders/qa/abstract"

type Spider interface {
	Run() error
	Update() error
	GetAllQA() []abstract.ArticleQA
}
