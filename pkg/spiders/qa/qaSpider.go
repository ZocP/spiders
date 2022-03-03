package qa

import (
	"go.uber.org/zap"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
	"qa_spider/pkg/spiders/qa/dynamics"
)

type QaSpider struct {
	storage []abstract.ArticleQA
	config  *config.Config
	log     *zap.Logger
}

func (q *QaSpider) Run() error {
	ids, err := dynamics.GetDynamicsIDs(q.log, q.config)
	if err != nil {
		return err
	}
	q.storage = dynamics.GetArticle(q.log, ids)
	return nil
}

func (q *QaSpider) Update() error {
	panic("implement me")
}

func (q *QaSpider) GetAllQA() []abstract.ArticleQA {
	return q.storage
}
