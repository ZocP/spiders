package qa

import (
	"fmt"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
	"qa_spider/pkg/spiders/qa/dynamics"
	"qa_spider/pkg/spiders/qa/writer"
	"regexp"
)

const (
	INAME = "QASPIDER"
)

type QaSpider struct {
	IName   string
	writer  writer.Writer
	storage []abstract.ArticleQA
	config  *config.Config
	log     *zap.Logger
	cron    *cron.Cron
}

func (q *QaSpider) GetName() string {
	return INAME
}

func (q *QaSpider) Run() error {
	q.log.Info("QA spider running，羊驼k48")
	if q.config.Internal.QASpider.Writer.Type == config.LocalTxt {
		//f, err := os.Open(q.config.Internal.QASpider.Writer.LocalTxt.Path + QAFileName)
		if f := q.writer.ReadArticleQA(q.config.Internal.QASpider.Writer.LocalTxt.Path); f != nil {
			q.log.Info("loaded from local storage")
			q.storage = f
			return nil
		}
		q.log.Info("file not exist, generating new QA")
		ids, err := dynamics.GetDynamicsIDs(q.log, q.config)
		if err != nil {
			return err
		}
		q.storage = dynamics.GetArticle(q.log, ids)
		if err := q.writer.WriteArticleQA(q.storage); err != nil {
			q.log.Error("Error writing QA", zap.Error(err))
		}
	}
	if q.config.Internal.QASpider.AutoUpdate {
		err := q.timeUpdater("2")
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *QaSpider) Update() error {
	q.log.Info("spider updating")
	cv := regexp.MustCompile("[0-9].*")
	if len(q.storage) == 0 {
		ids, err := dynamics.GetDynamicsIDs(q.log, q.config)
		if err != nil {
			q.log.Info("error happen while updating", zap.Error(err))
			return err
		}
		q.storage = dynamics.GetArticle(q.log, ids)
	}
	scv := cv.FindStringSubmatch(q.storage[0].Link)[0]
	q.log.Info("latest QA in storage", zap.String("cv", scv))
	ids, err := dynamics.GetDynamicsIDs(q.log, q.config, scv)
	if err != nil {
		q.log.Info("error happen while updating", zap.Error(err))
		return err
	}
	appender := dynamics.GetArticle(q.log, ids)
	if len(appender) == 0 {
		return fmt.Errorf("nothing can be updated")
	}
	q.storage = append(appender, q.storage...)
	if err := q.writer.WriteArticleQA(q.storage); err != nil {
		return err
	}
	return nil
}

func (q *QaSpider) GetAllQA() []abstract.ArticleQA {
	return q.storage
}

func (q *QaSpider) Reload() error {
	if f := q.writer.ReadArticleQA(q.config.Internal.QASpider.Writer.LocalTxt.Path); f != nil {
		q.log.Info("reload from file", zap.Int("total articles", len(f)))
		q.storage = f
		return nil
	}
	return fmt.Errorf("unable to reload")
}
func InitDefaultSpider(writer writer.Writer, config *config.Config, log *zap.Logger) Spider {
	return &QaSpider{
		writer:  writer,
		storage: make([]abstract.ArticleQA, 0),
		config:  config,
		log:     log,
	}
}

func (q *QaSpider) timeUpdater(date string) error {
	filter := regexp.MustCompile("[1-7]")
	if filter.FindStringSubmatch(date)[0] == "" {
		return fmt.Errorf("invalid date input")
	}
	c := cron.New()
	q.cron = c
	if err := c.AddFunc("0 0 0 0 0 WED", func() {
		if err := q.Update(); err != nil {
			q.log.Info("error ")
		}
	}); err != nil {
		q.log.Error("adding func", zap.Error(err))
	}
	return nil
}
