package writer

import (
	"qa_spider/pkg/spiders/qa/abstract"
)

type Writer interface {
	WriteArticleQA(articles []abstract.ArticleQA) error
	ReadArticleQA(path string) []abstract.ArticleQA
}
