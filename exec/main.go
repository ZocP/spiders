package main

import (
	"go.uber.org/zap"
	"log"
	"qa_spider/config"
	"qa_spider/pkg/console"
	"qa_spider/pkg/spiders/qa"
	"qa_spider/pkg/spiders/qa/writer"
	"qa_spider/server"
)

func main() {
	if err := InitDependencies().Run(); err != nil {
		log.Fatal("fatal while starting")
	}
}

func InitDependencies() server.Server {

	log, _ := zap.NewDevelopment()
	c := config.InitConfig(log)

	writer := writer.InitDefaultWriter(log, c)
	spider := qa.InitDefaultSpider(writer, c, log)

	server := server.InitHTTPServer(c, log, spider)

	listener := console.InitListener(log, spider, server)
	listener.Run()

	return server
}
