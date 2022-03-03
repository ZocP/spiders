package writer

import (
	"qa_spider/pkg/spiders/qa/abstract"
)

type Writer interface {
	WriteArticleQA(articles []abstract.ArticleQA, args ...interface{}) error
	ReadArticleQA(path string) []abstract.ArticleQA
}