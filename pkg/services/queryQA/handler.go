package queryQA

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"qa_spider/pkg/internalsvc/spiders/qa"
	"qa_spider/pkg/services"
	"qa_spider/server/content"
	"strings"
	"time"
	"unicode/utf8"
)

func QueryQA(ctn *content.Content) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		kw := ctx.Query("key")
		kw, err := url.QueryUnescape(kw)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, services.ErrorResponse(err))
		}
		length := utf8.RuneCountInString(kw)
		if length < ctn.Config.Services.QueryQA.Shortest || length > ctn.Config.Services.QueryQA.Longest {
			ctx.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("too short or too long keyword")))
			return
		}
		ctn.Debug("request key", zap.String("key", kw))
		ctn.Debug("request key", zap.Int("len", length))
		if kw == "小伙伴你好" {
			ctx.JSON(http.StatusOK, services.SuccessResponse("YBB"))
			return
		}
		result := findMatchesWithTime(kw, ctn)
		ctn.Debug("matches found", zap.Int("result found", len(result)))
		if result == nil || len(result) == 0 {
			ctx.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("no match QA found")))
			return
		}
		if ctn.Config.Services.QueryQA.MaxArticleResults > 0 && len(result) > ctn.Config.Services.QueryQA.MaxArticleResults {
			ctx.JSON(http.StatusOK, services.ErrorResponse(fmt.Errorf("too many results, use a longer keyword")))
			return
		}

		ctx.JSON(http.StatusOK, services.SuccessResponse(result))
	}
}

func findMatchesWithTime(find string, ctn *content.Content) []Result {
	var r []Result
	begin := time.Now()
	articles := ctn.Data[0].(qa.Spider).GetAllQA()
	for _, v := range articles {
		result := Result{
			Title: v.Title,
			Link:  v.Link,
			QA:    make([]PairQA, 0),
		}
		for _, v2 := range v.QA {
			if strings.Index(v2.Q, find) > 0 && len(result.QA) < ctn.Config.Services.QueryQA.MaxQAPair {
				result.QA = append(result.QA, PairQA{
					Q: v2.Q,
					A: v2.A,
				})
				continue
			}
			if strings.Index(v2.A, find) > 0 && len(result.QA) < ctn.Config.Services.QueryQA.MaxQAPair {
				result.QA = append(result.QA, PairQA{
					Q: v2.Q,
					A: v2.A,
				})
				continue
			}
		}
		if len(result.QA) == 0 {
			continue
		}
		r = append(r, result)
	}
	end := time.Now()
	ctn.Debug("regular search time used:", zap.Int64("in milliseconds", end.UnixMilli()-begin.UnixMilli()))
	return r
}
