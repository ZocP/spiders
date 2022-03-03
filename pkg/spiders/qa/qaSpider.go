package qa

import (
	"go.uber.org/zap"
	"qa_spider/config"
)

type QaSpider struct {
	config *config.Config
	log    *zap.Logger
}

func (q *QaSpider) Run() error {
	//dynamicIds := dynamics.GetDynamicsIDs(q.log, config.Config{})
	return nil
}

func (q *QaSpider) Update() error {
	panic("implement me")
}
