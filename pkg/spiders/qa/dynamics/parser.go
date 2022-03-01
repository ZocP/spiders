package dynamics

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
	"regexp"
	"strings"
)

func GetDynamicsIDs(log *zap.Logger, config *config.Config) []string {
	pn := 1
	ps := 30
	var dynamics []string
	for {
		body, err := getQADynamicAPI(pn, ps, log)
		if err != nil {
			log.Error("requesting QA dynamic ", zap.Error(err))
			//TODO: add retry
			//switch err{
			//
			//}
		}
		val := gjson.Parse(body)
		if val.Get("code").String() != "0" {
			log.Error("parsing json", zap.String("error", val.Get("message").String()))
			return nil
		}
		if val.Get("data.cards").Value() == nil {
			log.Info("end of cards")
			log.Info("dynamics: ", zap.Any("all", dynamics))
			return dynamics
		}
		val.Get("data.cards").ForEach(func(key, value gjson.Result) bool {
			match := regexp.MustCompile(`制作委员会的每周QA\s\d+\.\d+`)
			result := match.FindAllStringSubmatch(value.Get("card").String(), -1)
			if result == nil {
				return true
			}
			log.Info("found matches: ", zap.Any("result", result))
			s := value.Get("desc.rid").String()

			//避免叔叔给的过长的cv号，叔叔真是4🐎了哈哈哈

			if len(s) > 12 {
				return true
			}

			dynamics = append(dynamics, s+":"+result[0][0])
			return true
		})
		pn++
	}
}

func GetArticle(log *zap.Logger, dIDs []string) []abstract.ArticleQA {
	var result []abstract.ArticleQA
	for _, val := range dIDs {
		cv := val[:strings.Index(val, ":")]
		title := val[strings.Index(val, ":")+1:]
		log.Info("cv", zap.String("cv", cv))
		log.Info("title", zap.String("title", title))
		current := abstract.ArticleQA{
			Link:  "https://www.bilibili.com/read/cv" + cv,
			Title: title,
			QA:    make([]abstract.PairQA, 0),
			Mark:  "END HERE",
		}
		body, err := getArticle(cv, log)
		if err != nil {
			//TODO retry
		}
		if get := getQAPairs(filter(body, log), log); get == nil {
			continue
		} else {
			current.QA = get
		}
		result = append(result, current)
	}
	return nil
}

func filter(reader io.ReadCloser, log *zap.Logger) string {
	defer reader.Close()
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Error("err parsing document", zap.Error(err))
		return ""
	}
	node := document.Find("div.read-article-holder").Text()
	return node
}

func getQAPairs(raw string, log *zap.Logger) []abstract.PairQA {
	var all []abstract.PairQA
	if raw == "" {
		log.Info("empty article")
		return nil
	}

	i := strings.Index(raw, "Q：")
	nraw := raw[i:]
	result := strings.Split(nraw, "Q：")
	for _, v := range result {
		if v == "" {
			continue
		}
		QA := strings.Split(v, "A：")

		new := abstract.PairQA{
			Q: make([]string, 0),
			A: QA[1],
		}
		new.Q = append(new.Q, QA[0])
		all = append(all, new)
	}
	log.Info("split", zap.Any("here", all))
	return all
}
