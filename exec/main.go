package main

import (
	"go.uber.org/zap"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa"
	"qa_spider/pkg/spiders/qa/writer"
)

func main() {
	log, _ := zap.NewDevelopment()
	c := config.InitConfig(log)
	writer := writer.InitDefaultWriter(log, c)
	//ids, _ := dynamics.GetDynamicsIDs(log, c)
	//articles := dynamics.GetArticle(log, ids)
	spider := qa.InitDefaultSpider(writer, c, log)
	spider.Run()
	if err := spider.Update(); err != nil {
		log.Info("update", zap.Error(err))
	}
	if err := writer.WriteArticleQA(spider.GetAllQA()); err != nil {
		log.Error("err while writing article", zap.Error(err))
	}
}
