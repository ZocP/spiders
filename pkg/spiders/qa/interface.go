package qa

import (
	"qa_spider/pkg"
	"qa_spider/pkg/spiders/qa/abstract"
)

type Spider interface {
	pkg.Internal
	Run() error
	Update() error
	GetAllQA() []abstract.ArticleQA
	Reload() error
}
