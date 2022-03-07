package console

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"qa_spider/config"
	"qa_spider/pkg/internalsvc/spiders/qa"
	"qa_spider/pkg/internalsvc/spiders/qa/writer"
	"qa_spider/server"
)

type Listener struct {
	*zap.Logger
	config *config.Config
	spider qa.Spider
	server server.Server
	writer writer.Writer
}

func (l *Listener) Run() {
	go func() {
		for {
			var operator string
			var function string
			fmt.Scanf("%s %s", &operator, &function)
			switch operator {
			case "spider":
				switch function {
				case "update":
					if err := l.spider.Update(); err != nil {
						l.Info("updating", zap.Error(err))
					}
				case "reload":
					if err := l.spider.Reload(); err != nil {
						l.Info("reloading", zap.Error(err))
					}
				case "inspect_all_titles":
					qa := l.spider.GetAllQA()
					for _, v := range qa {
						l.Info("inspecting", zap.String("title", v.Title))
					}
				case "write_test_file":
					l.spider.GetAllQA()
				default:
					l.Info("didn't find this method, available functions: update; reload; inspect_all_titles")
				}
			case "server":
				switch function {
				case "stop":
					os.Exit(1)
				default:
					l.Info("did not find this method")
				}

			default:
				l.Info("didn't find this operator")
			}
		}
	}()
}

func InitListener(log *zap.Logger, spider qa.Spider, server server.Server) *Listener {
	return &Listener{
		Logger: log,
		spider: spider,
		server: server,
	}
}
