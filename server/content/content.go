package content

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"qa_spider/config"
	"strconv"
)

type Content struct {
	Log    *zap.Logger
	Config *config.Config
	Db     *gorm.DB
	Data   []interface{}
}

func InitContent(config *config.Config, log *zap.Logger, service config.ServiceID) *Content {
	dsn := getDSN(config, service, log)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("connecting to database: ", zap.Error(err))
	}
	return &Content{
		Config: config,
		Db:     db,
		Log:    log,
		Data:   make([]interface{}, 0),
	}
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
