package console

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"qa_spider/pkg/spiders/qa"
	"qa_spider/server"
)

type Listener struct {
	*zap.Logger
	spider qa.Spider
	server server.Server
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
				default:
					l.Info("didn't find this method")
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
