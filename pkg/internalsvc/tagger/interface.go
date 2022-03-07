package tagger

import (
	"qa_spider/pkg"
	"qa_spider/pkg/internalsvc/spiders/qa/abstract"
	"qa_spider/pkg/internalsvc/tagger/tag"
)

type Tagger interface {
	pkg.Internal
	UpdateTagWithQAs(tags []tag.Tag, qas []abstract.ArticleQA) error
	GetTagWithQAs(tag string) ([]abstract.ArticleQA, error)
}
