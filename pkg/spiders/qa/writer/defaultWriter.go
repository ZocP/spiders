package writer

import (
	"bufio"
	"go.uber.org/zap"
	"os"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
	"regexp"
)

type DefaultWriter struct {
	*zap.Logger
	config *config.Config
}

func (d *DefaultWriter) WriteArticleQA(articles []abstract.ArticleQA, args ...interface{}) error {
	var file *os.File
	if args == nil {
		err := os.MkdirAll(d.config.Internal.QASpider.Writer.LocalTxt.Path, 0777)
		if err != nil {
			return err
		}
		f, err := os.Create(d.config.Internal.QASpider.Writer.LocalTxt.Path + "QA.txt")
		if err != nil {
			return err
		}
		file = f
		defer file.Close()
	} else {
		err := os.MkdirAll(args[0].(string), 0777)
		if err != nil {
			return err
		}
		f, err := os.Create(args[0].(string) + "QA.txt")
		if err != nil {
			return err
		}
		file = f
		defer file.Close()
	}
	w := bufio.NewWriter(file)
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
	f, err := os.Open(path + "QA.txt")
	if err != nil {
		d.Error("opening file", zap.Error(err))
	}
	reader := bufio.NewScanner(f)
	title := regexp.MustCompile(`制作委员会的每周QA\s\d+\.\d+`)
	articles := make([]abstract.ArticleQA, 0)
	for reader.Scan() {
		if title.FindAllStringSubmatchIndex(reader.Text(), -1) != nil {
			stitle := reader.Text()
			reader.Scan()
			article := abstract.ArticleQA{
				Link:  reader.Text(),
				Title: stitle,
				QA:    make([]abstract.PairQA, 0),
				Mark:  "END HERE",
			}
			for reader.Scan() {
				if len(reader.Text()) < 2 {
					break
				}
				Q := make([]string, 1)
				A := ""
				Q[0] = ""
				//如果这行是Q
				if reader.Text()[0:1] == "Q" {
					Q[0] += reader.Text()
					for reader.Scan() {
						if len(reader.Text()) < 1 {
							break
						}
						if reader.Text()[0:1] != "E" && reader.Text()[0:1] != "A" {
							Q[0] += reader.Text()
						} else {
							break
						}
					}
				}
				//如果这行是A
				if reader.Text()[0:1] == "A" {
					A += reader.Text()
					article.QA = append(article.QA, abstract.PairQA{
						Q: Q,
						A: A,
					})
					Q = make([]string, 0)
					A = ""
					//for reader.Scan() {
					//	if len(reader.Text()) < 1 {
					//		break
					//	}
					//	if reader.Text()[0:1] != "E" && reader.Text()[0:1] != "Q" {
					//		A += reader.Text()
					//		break
					//	} else {
					//		article.QA = append(article.QA, abstract.PairQA{
					//			Q: Q,
					//			A: A,
					//		})
					//		Q = make([]string, 0)
					//		A = ""
					//		break
					//	}
					//}
				}
				if len(reader.Text()) < 1 {
					reader.Scan()
				}

				//如果这行是结束
				if reader.Text()[0:1] == "E" {
					break
				}
			}
			articles = append(articles, article)
		}
	}
	return articles
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
