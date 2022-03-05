package content

import (
	"go.uber.org/zap"
	"qa_spider/config"
	"strconv"
)

type Content struct {
	*zap.Logger
	Config *config.Config
	Data   []interface{}
}

func InitContent(config *config.Config, log *zap.Logger, data ...interface{}) *Content {
	//dsn := getDSN(config, log)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.Error("connecting to database: ", zap.Error(err))
	//}
	r := &Content{
		Config: config,
		Logger: log,
		Data:   make([]interface{}, 0),
	}
	r.Data = data
	return r
}

func getDSN(config *config.Config, service config.ServiceID, log *zap.Logger) string {
	DB := config.GetServiceDB(service)
	un := DB.UserName
	pc := DB.Password
	prtc := DB.Protocol
	url := DB.URL
	dn := DB.DBName
	r := un + ":" + pc + "@" + prtc + "(" + url + ")/" + dn
	log.Info("service: " + strconv.Itoa(int(service)) + " service dsn: " + r)
	return r
}
