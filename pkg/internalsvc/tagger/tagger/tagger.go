package tagger

import (
	"go.uber.org/zap"
	"qa_spider/pkg/internalsvc/spiders/qa/abstract"
	"qa_spider/pkg/internalsvc/tagger"
	"qa_spider/pkg/internalsvc/tagger/tag"
)

const (
	IName = "TAGGER"
)

type DefaultTagger struct {
	*zap.Logger
}

func (d *DefaultTagger) Run() error {

}

func (d *DefaultTagger) GetName() string {
	return IName
}

func (d *DefaultTagger) UpdateTagWithQAs(tags []tag.Tag, qas []abstract.ArticleQA) error {
	panic("implement me")
}

func (d *DefaultTagger) GetTagWithQAs(tag string) ([]abstract.ArticleQA, error) {
	panic("implement me")
}

func InitDefaultTagger() tagger.Tagger {
}
