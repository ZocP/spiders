package dynamics

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"qa_spider/config"
)

func getDynamicsIDs(log *zap.Logger, config config.Config) []string{
	pn := 0
	ps := 30
	var dynamics []int
	for{
		body, err := getQADynamicAPI(pn, ps, log)
		if err != nil{
			log.Error("requesting QA dynamic ", zap.Error(err))
			//TODO: add retry
			//switch err{
			//
			//}
		}
		val := gjson.Parse(body)
		val.Get("data.cards").ForEach(func(key, value gjson.Result) bool {

		})
	}

}