package queryWithOption

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"go.uber.org/zap"
	"net/http"
	"qa_spider/pkg/internalsvc/spiders/qa"
	"qa_spider/pkg/services"
	"qa_spider/pkg/services/queryQA"
	"qa_spider/server/content"
	"regexp"
	"time"
	"unicode/utf8"
)

func QueryWithOption(ctn *content.Content) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestBody
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("error format of JSON")))
			return
		}

		length := utf8.RuneCountInString(req.Keyword)
		if length < ctn.Config.Services.QueryQA.Shortest || length > ctn.Config.Services.QueryQA.Longest {
			c.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("too short or too long keyword")))
			return
		}
		ctn.Debug("request key", zap.String("key, option", req.Keyword+","+req.Option))
		ctn.Debug("request keyword", zap.Int("len", length))
		switch req.Option {
		case REGEX:
			result, err := queryWithRegex(ctn, req.Keyword)
			ctn.Debug("regex matches found", zap.Int("result found", len(result)))
			if err != nil {
				c.JSON(http.StatusOK, services.ErrorResponse(err))
				return
			}
			if len(result) == 0 {
				c.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("no match QA found")))
				return
			}
			c.JSON(http.StatusOK, services.SuccessResponse(result))
			return
		case FUZZY:
			result := fuzzyQuery(ctn, req.Keyword)
			ctn.Debug("fuzzy matches found", zap.Int("result found", len(result)))
			if result == nil || len(result) == 0 {

				c.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("no match QA found")))
				return
			}
			c.JSON(http.StatusOK, services.SuccessResponse(fuzzyQuery(ctn, req.Keyword)))
			return
		default:
			c.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("option %s not support", req.Option)))
			return
		}
	}
}

func queryWithRegex(ctn *content.Content, key string) (res []queryQA.Result, err error) {
	var resp []queryQA.Result
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("parsing regex err")
		}
	}()
	compiler := regexp.MustCompile(key)
	if compiler == nil {
		return nil, fmt.Errorf("illegal regex")
	}

	articles := ctn.Data[0].(qa.Spider).GetAllQA()
	begin := time.Now()
	for _, v := range articles {
		result := queryQA.Result{
			Title: v.Title,
			Link:  v.Link,
			QA:    make([]queryQA.PairQA, 0),
		}
		for _, v2 := range v.QA {
			if compiler.MatchString(v2.Q) && len(result.QA) < ctn.Config.Services.QueryQA.MaxQAPair {
				result.QA = append(result.QA, queryQA.PairQA{
					Q: v2.Q,
					A: v2.A,
				})
				continue
			}
			if compiler.MatchString(v2.A) && len(result.QA) < ctn.Config.Services.QueryQA.MaxQAPair {
				result.QA = append(result.QA, queryQA.PairQA{
					Q: v2.A,
					A: v2.Q,
				})
				continue
			}
		}
		if len(result.QA) == 0 {
			continue
		}
		resp = append(resp, result)
	}
	end := time.Now()
	ctn.Debug("regex query time used:", zap.Int64("in milliseconds", end.UnixMilli()-begin.UnixMilli()))
	return resp, err
}

func fuzzyQuery(ctn *content.Content, key string) []queryQA.Result {
	var resp []queryQA.Result
	articles := ctn.Data[0].(qa.Spider).GetAllQA()
	begin := time.Now()
	for _, v := range articles {
		result := queryQA.Result{
			Title: v.Title,
			Link:  v.Link,
			QA:    make([]queryQA.PairQA, 0),
		}
		for _, v2 := range v.QA {
			if fuzzy.Match(key, v2.Q) && len(result.QA) < ctn.Config.Services.QueryQA.MaxQAPair {
				result.QA = append(result.QA, queryQA.PairQA{
					Q: v2.Q,
					A: v2.A,
				})
				continue
			}
			if fuzzy.Match(key, v2.A) && len(result.QA) < ctn.Config.Services.QueryQA.MaxQAPair {
				result.QA = append(result.QA, queryQA.PairQA{
					Q: v2.A,
					A: v2.Q,
				})
				continue
			}
		}
		if len(result.QA) == 0 {
			continue
		}
		resp = append(resp, result)
	}
	end := time.Now()
	ctn.Debug("fuzzy query time used:", zap.Int64("in milliseconds", end.UnixMilli()-begin.UnixMilli()))
	return resp
}
