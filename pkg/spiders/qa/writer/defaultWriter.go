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

func (d *DefaultWriter) WriteArticleQA(articles []*abstract.ArticleQA, args ...interface{}) error {
	var file *os.File
	if args == nil {
		err := os.MkdirAll(d.config.Internal.QASpider.Writer.LocalTxt.Path, 0777)
		if err != nil {
			return err
		}
		f, err := os.Create(d.config.Internal.QASpider.Writer.LocalTxt.Path + QAFileName)
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
		f, err := os.Create(args[0].(string) + QAFileName)
		if err != nil {
			return err
		}
		file = f
		defer file.Close()
	}
	w := bufio.NewWriter(file)
	defer w.Flush()

	for _, v := range articles {
		d.Debug("writing article", zap.String("article", v.Title))
		_, err := w.WriteString(
			v.Title + "\n" +
				v.Link + "\n" +
				pairQAToString(d.Logger, v.QA) +
				v.Mark + "\n\n")
		if err != nil {
			d.Error("err while writing", zap.Error(err))
		}
	}
	d.Info("finished writing articles", zap.String("path", d.config.Internal.QASpider.Writer.LocalTxt.Path))
	return nil
}

func (d *DefaultWriter) ReadArticleQA(path string) []*abstract.ArticleQA {
	f, err := os.Open(path + "QA.txt")
	if err != nil {
		d.Info("opening file", zap.Error(err))
		return nil
	}
	reader := bufio.NewScanner(f)
	title := regexp.MustCompile(`制作委员会的每周QA\s\d+\.\d+`)
	cv := regexp.MustCompile("[0-9].*")
	articles := make([]*abstract.ArticleQA, 0)
	for reader.Scan() {
		if title.FindAllStringSubmatchIndex(reader.Text(), -1) != nil {
			stitle := reader.Text()
			reader.Scan()
			article := &abstract.ArticleQA{
				CV:    cv.FindStringSubmatch(reader.Text())[0],
				Link:  reader.Text(),
				Title: stitle,
				QA:    make([]abstract.PairQA, 0),
				Mark:  "END HERE",
			}
			for reader.Scan() {
				if len(reader.Text()) < 2 {
					break
				}
				Q := ""
				A := ""
				//如果这行是Q
				if reader.Text()[0:1] == "Q" {
					Q += reader.Text()
					for reader.Scan() {
						if len(reader.Text()) < 1 {
							break
						}
						if reader.Text()[0:1] != "E" && reader.Text()[0:1] != "A" {
							Q += reader.Text()
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
					Q = ""
					A = ""
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

//TODO: 添加append
func (d *DefaultWriter) AppendQA(qa abstract.ArticleQA) error {
	return nil
}

func pairQAToString(log *zap.Logger, qa []abstract.PairQA) string {
	result := ""
	for _, pair := range qa {

		result += pair.Q + "\n"
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
