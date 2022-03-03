package queryQA

import (
	"github.com/gin-gonic/gin"
	"log"
	"qa_spider/pkg/spiders/qa/abstract"
	"qa_spider/server/content"
	"strings"
	"time"
)

func QueryQA(ctx content.Content) gin.HandlerFunc {
	panic("implement me")
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
					log.Println("found at: " + v3)
				}
			}
			if strings.Index(v2.A, find) > 0 {
				log.Println("found at: ", v2.A)
			}
		}
	}
	end := time.Now()
	log.Println("time used:", end.UnixMilli()-begin.UnixMilli())
}
