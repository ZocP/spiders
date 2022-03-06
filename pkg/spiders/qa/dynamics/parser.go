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

//arguments 里面可以塞一个终止值，如果读取到立刻停止
func GetDynamicsIDs(log *zap.Logger, config *config.Config, args ...interface{}) ([]string, error) {
	pn := 1
	ps := 30
	var dynamics []string
	con := true
	if args != nil && args[0] != nil {
		log.Debug("input info with arguments")
		for con {
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
				return nil, err
			}
			if val.Get("data.cards").Value() == nil {
				log.Info("end of cards")
				log.Info("dynamics: ", zap.Any("all", dynamics))
				return dynamics, nil
			}
			val.Get("data.cards").ForEach(func(key, value gjson.Result) bool {
				match := regexp.MustCompile(`制作委员会的每周QA\s\d+\.\d+`)
				result := match.FindAllStringSubmatch(value.Get("card").String(), -1)
				if result == nil {
					return true
				}
				log.Debug("found matches: ", zap.Any("result", result))
				s := value.Get("desc.rid").String()
				//避免叔叔给的过长的cv号，叔叔真是4🐎了哈哈哈
				if s == args[0].(string) {
					con = false
					return false
				}
				if len(s) > 12 {
					return true
				}

				dynamics = append(dynamics, s+":"+result[0][0])
				return true
			})
			pn++
		}
	}
	if con == false {
		return dynamics, nil
	}

	log.Debug("input info without arguments")
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
			return nil, err
		}
		if val.Get("data.cards").Value() == nil {
			log.Info("end of cards")
			log.Info("dynamics: ", zap.Any("all", dynamics))
			return dynamics, nil
		}
		val.Get("data.cards").ForEach(func(key, value gjson.Result) bool {
			match := regexp.MustCompile(`制作委员会的每周QA\s\d+\.\d+`)
			result := match.FindAllStringSubmatch(value.Get("card").String(), -1)
			if result == nil {
				return true
			}
			log.Debug("found matches: ", zap.Any("result", result))
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

func GetArticle(log *zap.Logger, dIDs []string) []*abstract.ArticleQA {
	var result []*abstract.ArticleQA
	for _, val := range dIDs {
		cv := val[:strings.Index(val, ":")]
		title := val[strings.Index(val, ":")+1:]
		log.Debug("getting article", zap.String("article cv", cv))
		log.Debug("getting article", zap.String("article title", title))
		current := abstract.ArticleQA{
			CV:    cv,
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
		result = append(result, &current)
	}
	return result
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
	about := regexp.MustCompile("[^\u4e00-\uFFFFF\\w\uFF00-\uFFFF]{1,5}关于【?.{0,7}】?[^\u4ef00-\uFFFF\\w\uFF00-\uFFFF]{1,5}")
	raw = about.ReplaceAllString(raw, "")
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
		contain := regexp.MustCompile("A.?：.*")
		QA := contain.FindAllStringSubmatch(v, -1)

		index := contain.FindAllStringSubmatchIndex(v, -1)
		Q := v[:index[0][0]]
		new := abstract.PairQA{
			Q: "",
			A: QA[0][0],
		}
		new.Q = "Q：" + strings.Replace(Q, "\n", "", -1)
		new.Q = strings.Replace(new.Q, " ", "", -1)
		new.A = strings.Replace(new.A, " ", "", -1)
		all = append(all, new)
	}
	return all
}
