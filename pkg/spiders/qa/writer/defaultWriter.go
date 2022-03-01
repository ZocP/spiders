package writer

import (
	"bufio"
	"go.uber.org/zap"
	"os"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
)

type DefaultWriter struct{
	zap.Logger
	config config.Config
}

func (d *DefaultWriter) WriteArticleQA(articles []abstract.ArticleQA) error {

	f ,err := os.Create(d.config.Services.Writer.Path + "QA.txt")
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, v := range articles{
		w.WriteString(
			v.Title +
				"\n"+
				)

	}
	check(err)
}

func (d *DefaultWriter) ReadArticleQA(path string) []abstract.ArticleQA {
	panic("implement me")
}
