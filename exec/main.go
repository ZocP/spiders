package main

import (
	"go.uber.org/zap"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/dynamics"
	"qa_spider/pkg/spiders/qa/writer"
)

func main() {
	log, _ := zap.NewDevelopment()
	c := config.InitConfig(log)
	writer := writer.InitDefaultWriter(log, c)
	ids := dynamics.GetDynamicsIDs(log, c)
	articles := dynamics.GetArticle(log, ids)
	if err := writer.WriteArticleQA(articles); err != nil {
		log.Error("err while writing article", zap.Error(err))
	}

}
