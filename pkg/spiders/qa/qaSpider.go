package qa

import (
	"go.uber.org/zap"
	"qa_spider/config"
)

type QaSpider struct {
	config *config.Config
	log *zap.Logger
}

func (q *QaSpider) Run() error {
	dynamics
}

func (q *QaSpider) Update() error {
	panic("implement me")
}

