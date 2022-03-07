package tag

import "qa_spider/pkg/internalsvc/spiders/qa/abstract"

type Tag struct {
	Tag     string
	KeyWord []string
}

type TagWithArticle struct {
	Tag string
	QAs []abstract.ArticleQA
}
