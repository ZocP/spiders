package main

import (
	"go.uber.org/zap"
	"log"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
	"qa_spider/pkg/spiders/qa/writer"
	"strings"
	"time"
)

func main() {
	log, _ := zap.NewDevelopment()
	c := config.InitConfig(log)
	writer := writer.InitDefaultWriter(log, c)
	//ids := dynamics.GetDynamicsIDs(log, c)
	//articles := dynamics.GetArticle(log, ids)
	//if err := writer.WriteArticleQA(articles); err != nil {
	//	log.Error("err while writing article", zap.Error(err))
	//}
	article2 := writer.ReadArticleQA("./files/spider/")
	writer.WriteArticleQA(article2, "./files/test/")
}

func findMatchesWithTime(find string, articles []abstract.ArticleQA) {
	begin := time.Now()
	for _, v := range articles {
		if strings.Index(v.Title, find) > 0 {
			log.Println("found at" + v.Title)
		}
		for _, v2 := range v.QA {
			for _, v3 := range v2.Q {
				if strings.Index(v3, find) > 0 {
					log.Println("found at" + v3)
				}
			}
			if strings.Index(v2.A, find) > 0 {
				log.Println("found at", v2.A)
			}
		}
	}
	end := time.Now()
	log.Println("time used:", end.UnixMicro()-begin.UnixMicro())
}
