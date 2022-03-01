package dynamics

import (
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"qa_spider/config"
	"qa_spider/pkg/spiders/qa/abstract"
	"regexp"
)

func GetDynamicsIDs(log *zap.Logger, config config.Config) []string {
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
				log.Info("skipping")
				return true
			}
			log.Info("found matches: ", zap.Any("result", result))
			s := value.Get("desc.rid").String()

			//避免叔叔给的过长的cv号，叔叔真是4🐎了哈哈哈

			if len(s) > 12 {
				return true
			}

			dynamics = append(dynamics, s)
			return true
		})
		pn++
	}
}

func GetArticle(log *zap.Logger, dIDs []string) []abstract.ArticleQA {
	var result []abstract.ArticleQA
	for _, val := range dIDs {
		current := abstract.ArticleQA{
			Link:  "https://www.bilibili.com/read/cv" + val,
			Title: "",
			QA:    make([]abstract.PairQA, 0),
			Mark:  "END HERE",
		}
		body, err := getArticle(val, log)
		if err != nil {
			//TODO retry
		}
		pairs := getQAPairs(filter(body))
	}
}

func filter(raw string) string {

}

func getQAPairs(raw string) []abstract.PairQA {

}
