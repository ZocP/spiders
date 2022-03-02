package writer

import (
	"bufio"
	"go.uber.org/zap"
	"os"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
)

type DefaultWriter struct {
	*zap.Logger
	config *config.Config
}

func (d *DefaultWriter) WriteArticleQA(articles []abstract.ArticleQA) error {
	err := os.MkdirAll(d.config.Services.Writer.Path, 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(d.config.Services.Writer.Path + "QA.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, v := range articles {
		d.Info("writing article", zap.String("article", v.Title))
		_, err := w.WriteString(
			v.Title + "\n" +
				v.Link + "\n" +
				pairQAToString(d.Logger, v.QA) +
				v.Mark + "\n\n")
		if err != nil {
			d.Error("err while writing", zap.Error(err))
		}
	}
	return nil
}

func (d *DefaultWriter) ReadArticleQA(path string) []abstract.ArticleQA {
	panic("implement me")
}

func pairQAToString(log *zap.Logger, qa []abstract.PairQA) string {
	result := ""
	for _, pair := range qa {
		for _, v := range pair.Q {
			result += v + "\n"
		}
		result += pair.A + "\n"
	}
	return result
}

func InitDefaultWriter(logger *zap.Logger, config *config.Config) Writer {
	return &DefaultWriter{
		Logger: logger,
		config: config,
	}
}
